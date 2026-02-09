#include "grpcclient.h"
#include "api/proto/kernel.grpc.pb.h"
#include "api/proto/qemu.grpc.pb.h"
#include "api/proto/common.pb.h"
#include <grpcpp/grpcpp.h>
#include <QDebug>

class GrpcClient::Impl
{
public:
    std::shared_ptr<grpc::Channel> channel;
    std::unique_ptr<elmos::v1::KernelService::Stub> kernelStub;
    std::unique_ptr<elmos::v1::QEMUService::Stub> qemuStub;
};

GrpcClient::GrpcClient(const QString &serverAddress, QObject *parent)
    : QObject(parent)
    , impl_(std::make_unique<Impl>())
{
    impl_->channel = grpc::CreateChannel(serverAddress.toStdString(), grpc::InsecureChannelCredentials());
    impl_->kernelStub = elmos::v1::KernelService::NewStub(impl_->channel);
    impl_->qemuStub = elmos::v1::QEMUService::NewStub(impl_->channel);
}

GrpcClient::~GrpcClient() = default;

bool GrpcClient::isConnected() const
{
    if (!impl_->channel) {
        return false;
    }
    auto state = impl_->channel->GetState(false);
    return state == GRPC_CHANNEL_READY || state == GRPC_CHANNEL_IDLE;
}

QString GrpcClient::lastError() const
{
    return lastError_;
}

bool GrpcClient::buildKernel(const QStringList &targets, int jobs, const QString &arch, bool verbose)
{
    if (!impl_->kernelStub) {
        lastError_ = "gRPC stub not initialized";
        emit errorOccurred(lastError_);
        return false;
    }

    elmos::v1::BuildRequest request;
    for (const QString &target : targets) {
        request.add_targets(target.toStdString());
    }
    request.set_jobs(jobs);
    request.set_arch(arch.toStdString());
    request.set_verbose(verbose);

    grpc::ClientContext context;
    auto reader = impl_->kernelStub->Build(&context, request);

    elmos::v1::BuildProgress progress;
    while (reader->Read(&progress)) {
        switch (progress.event_case()) {
            case elmos::v1::BuildProgress::kStage: {
                const auto &stage = progress.stage();
                emit buildStageChanged(
                    QString::fromStdString(stage.name()),
                    stage.progress(),
                    stage.current_file(),
                    stage.total_files()
                );
                break;
            }
            case elmos::v1::BuildProgress::kLog: {
                const auto &log = progress.log();
                emit buildLogReceived(
                    static_cast<int>(log.level()),
                    QString::fromStdString(log.message()),
                    log.timestamp_ms()
                );
                break;
            }
            case elmos::v1::BuildProgress::kError: {
                const auto &error = progress.error();
                emit buildErrorReceived(
                    QString::fromStdString(error.message()),
                    QString::fromStdString(error.file()),
                    error.line()
                );
                break;
            }
            case elmos::v1::BuildProgress::kComplete: {
                const auto &complete = progress.complete();
                emit buildCompleted(
                    complete.success(),
                    complete.duration_ms(),
                    QString::fromStdString(complete.image_path())
                );
                break;
            }
            default:
                break;
        }
    }

    grpc::Status status = reader->Finish();
    if (!status.ok()) {
        lastError_ = QString("Build failed: %1").arg(QString::fromStdString(status.error_message()));
        emit errorOccurred(lastError_);
        return false;
    }

    return true;
}

bool GrpcClient::cloneKernel(const QString &version)
{
    if (!impl_->kernelStub) {
        lastError_ = "gRPC stub not initialized";
        emit errorOccurred(lastError_);
        return false;
    }

    elmos::v1::CloneRequest request;
    request.set_version(version.toStdString());

    grpc::ClientContext context;
    auto reader = impl_->kernelStub->Clone(&context, request);

    elmos::v1::CloneProgress progress;
    while (reader->Read(&progress)) {
        emit cloneProgress(progress.progress(), QString::fromStdString(progress.message()));
    }

    grpc::Status status = reader->Finish();
    if (!status.ok()) {
        lastError_ = QString("Clone failed: %1").arg(QString::fromStdString(status.error_message()));
        emit errorOccurred(lastError_);
        return false;
    }

    return true;
}

bool GrpcClient::configureKernel(const QString &configType)
{
    if (!impl_->kernelStub) {
        lastError_ = "gRPC stub not initialized";
        emit errorOccurred(lastError_);
        return false;
    }

    elmos::v1::ConfigureRequest request;
    request.set_config_type(configType.toStdString());

    grpc::ClientContext context;
    elmos::v1::ConfigureResponse response;
    grpc::Status status = impl_->kernelStub->Configure(&context, request, &response);

    if (!status.ok()) {
        lastError_ = QString("Configure failed: %1").arg(QString::fromStdString(status.error_message()));
        emit errorOccurred(lastError_);
        return false;
    }

    if (!response.success()) {
        lastError_ = QString::fromStdString(response.error_message());
        emit errorOccurred(lastError_);
        return false;
    }

    return true;
}

bool GrpcClient::cleanKernel(bool deepClean)
{
    if (!impl_->kernelStub) {
        lastError_ = "gRPC stub not initialized";
        emit errorOccurred(lastError_);
        return false;
    }

    elmos::v1::CleanRequest request;
    request.set_deep_clean(deepClean);

    grpc::ClientContext context;
    elmos::v1::CleanResponse response;
    grpc::Status status = impl_->kernelStub->Clean(&context, request, &response);

    if (!status.ok()) {
        lastError_ = QString("Clean failed: %1").arg(QString::fromStdString(status.error_message()));
        emit errorOccurred(lastError_);
        return false;
    }

    return response.success();
}

QStringList GrpcClient::listKernelVersions(int limit)
{
    QStringList versions;
    if (!impl_->kernelStub) {
        lastError_ = "gRPC stub not initialized";
        emit errorOccurred(lastError_);
        return versions;
    }

    elmos::v1::ListVersionsRequest request;
    request.set_limit(limit);

    grpc::ClientContext context;
    elmos::v1::ListVersionsResponse response;
    grpc::Status status = impl_->kernelStub->ListVersions(&context, request, &response);

    if (!status.ok()) {
        lastError_ = QString("ListVersions failed: %1").arg(QString::fromStdString(status.error_message()));
        emit errorOccurred(lastError_);
        return versions;
    }

    for (const auto &version : response.versions()) {
        versions.append(QString::fromStdString(version));
    }

    return versions;
}

bool GrpcClient::runQEMU(bool graphical, bool debug, int memoryMB, int cpus, const QStringList &extraArgs, const QString &kernelCmdline)
{
    if (!impl_->qemuStub) {
        lastError_ = "gRPC stub not initialized";
        emit errorOccurred(lastError_);
        return false;
    }

    elmos::v1::QEMURunRequest request;
    request.set_graphical(graphical);
    request.set_debug(debug);
    request.set_memory_mb(memoryMB);
    request.set_cpus(cpus);
    for (const QString &arg : extraArgs) {
        request.add_extra_args(arg.toStdString());
    }
    request.set_kernel_cmdline(kernelCmdline.toStdString());

    grpc::ClientContext context;
    auto reader = impl_->qemuStub->Run(&context, request);

    elmos::v1::QEMUOutput output;
    while (reader->Read(&output)) {
        switch (output.event_case()) {
            case elmos::v1::QEMUOutput::kStarted: {
                const auto &started = output.started();
                emit qemuStarted(
                    started.pid(),
                    QString::fromStdString(started.qemu_version()),
                    QString::fromStdString(started.command())
                );
                break;
            }
            case elmos::v1::QEMUOutput::kConsole: {
                const auto &console = output.console();
                QByteArray data(console.data().data(), console.data().size());
                emit qemuConsoleOutput(data);
                break;
            }
            case elmos::v1::QEMUOutput::kStopped: {
                const auto &stopped = output.stopped();
                emit qemuStopped(stopped.exit_code(), stopped.uptime_ms());
                break;
            }
            case elmos::v1::QEMUOutput::kError: {
                const auto &error = output.error();
                emit qemuError(QString::fromStdString(error.message()));
                break;
            }
            default:
                break;
        }
    }

    grpc::Status status = reader->Finish();
    if (!status.ok()) {
        lastError_ = QString("QEMU run failed: %1").arg(QString::fromStdString(status.error_message()));
        emit errorOccurred(lastError_);
        return false;
    }

    return true;
}

bool GrpcClient::stopQEMU()
{
    if (!impl_->qemuStub) {
        lastError_ = "gRPC stub not initialized";
        emit errorOccurred(lastError_);
        return false;
    }

    elmos::v1::QEMUStopRequest request;
    grpc::ClientContext context;
    elmos::v1::QEMUStopResponse response;
    grpc::Status status = impl_->qemuStub->Stop(&context, request, &response);

    if (!status.ok()) {
        lastError_ = QString("Stop QEMU failed: %1").arg(QString::fromStdString(status.error_message()));
        emit errorOccurred(lastError_);
        return false;
    }

    return response.success();
}

bool GrpcClient::sendQEMUInput(const QByteArray &data)
{
    if (!impl_->qemuStub) {
        lastError_ = "gRPC stub not initialized";
        emit errorOccurred(lastError_);
        return false;
    }

    elmos::v1::QEMUInputRequest request;
    request.set_data(data.constData(), data.size());

    grpc::ClientContext context;
    elmos::v1::QEMUInputResponse response;
    grpc::Status status = impl_->qemuStub->SendInput(&context, request, &response);

    if (!status.ok()) {
        lastError_ = QString("Send input failed: %1").arg(QString::fromStdString(status.error_message()));
        emit errorOccurred(lastError_);
        return false;
    }

    return response.success();
}
