#ifndef GRPCCLIENT_H
#define GRPCCLIENT_H

#include <QObject>
#include <QString>
#include <memory>

namespace grpc {
class Channel;
}

namespace elmos {
namespace v1 {
class KernelService;
}
}

class GrpcClient : public QObject
{
    Q_OBJECT

public:
    explicit GrpcClient(const QString &serverAddress = "unix:///tmp/elmos.sock", QObject *parent = nullptr);
    ~GrpcClient();

    bool isConnected() const;
    QString lastError() const;

    bool buildKernel(const QStringList &targets, int jobs, const QString &arch, bool verbose = false);
    bool cloneKernel(const QString &version);
    bool configureKernel(const QString &configType);
    bool cleanKernel(bool deepClean = false);
    QStringList listKernelVersions(int limit = 50);

    bool runQEMU(bool graphical, bool debug, int memoryMB, int cpus, const QStringList &extraArgs, const QString &kernelCmdline);
    bool stopQEMU();
    bool sendQEMUInput(const QByteArray &data);

signals:
    void buildStageChanged(const QString &stageName, int progress, int currentFile, int totalFiles);
    void buildLogReceived(int level, const QString &message, qint64 timestamp);
    void buildErrorReceived(const QString &message, const QString &file, int line);
    void buildCompleted(bool success, qint64 durationMs, const QString &imagePath);
    void cloneProgress(int progress, const QString &message);
    void qemuStarted(int pid, const QString &qemuVersion, const QString &command);
    void qemuConsoleOutput(const QByteArray &data);
    void qemuStopped(int exitCode, qint64 uptimeMs);
    void qemuError(const QString &error);
    void errorOccurred(const QString &error);

private:
    class Impl;
    std::unique_ptr<Impl> impl_;
    QString lastError_;
};

#endif
