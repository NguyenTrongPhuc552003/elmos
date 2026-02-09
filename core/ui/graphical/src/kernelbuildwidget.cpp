#include "kernelbuildwidget.h"
#include "grpcclient.h"
#include <QGroupBox>
#include <QDateTime>
#include <QThread>
#include <QScrollBar>

KernelBuildWidget::KernelBuildWidget(GrpcClient *grpcClient, QWidget *parent)
    : QWidget(parent)
    , grpcClient_(grpcClient)
    , isBuilding_(false)
{
    setupUI();
    connectSignals();
}

KernelBuildWidget::~KernelBuildWidget() = default;

void KernelBuildWidget::setupUI()
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);

    QGroupBox *configGroup = new QGroupBox("Build Configuration");
    QHBoxLayout *configLayout = new QHBoxLayout(configGroup);

    archCombo_ = new QComboBox();
    archCombo_->addItems({"arm", "arm64", "x86_64", "riscv", "mips"});
    archCombo_->setCurrentText("arm");
    configLayout->addWidget(new QLabel("Architecture:"));
    configLayout->addWidget(archCombo_);

    targetCombo_ = new QComboBox();
    targetCombo_->addItems({"Image", "zImage", "bzImage", "vmlinux", "modules", "dtbs"});
    targetCombo_->setCurrentText("Image");
    configLayout->addWidget(new QLabel("Target:"));
    configLayout->addWidget(targetCombo_);

    jobsSpinBox_ = new QSpinBox();
    jobsSpinBox_->setRange(0, 128);
    jobsSpinBox_->setValue(0);
    jobsSpinBox_->setSpecialValueText("Auto");
    configLayout->addWidget(new QLabel("Jobs:"));
    configLayout->addWidget(jobsSpinBox_);

    verboseCheckBox_ = new QCheckBox("Verbose");
    configLayout->addWidget(verboseCheckBox_);

    configLayout->addStretch();

    mainLayout->addWidget(configGroup);

    stageLabel_ = new QLabel("Ready");
    stageLabel_->setStyleSheet("font-weight: bold;");
    mainLayout->addWidget(stageLabel_);

    progressBar_ = new QProgressBar();
    progressBar_->setRange(0, 100);
    progressBar_->setValue(0);
    mainLayout->addWidget(progressBar_);

    outputView_ = new QTextEdit();
    outputView_->setReadOnly(true);
    outputView_->setFont(QFont("Monaco", 10));
    outputView_->setStyleSheet("background-color: #1e1e1e; color: #d4d4d4;");
    mainLayout->addWidget(outputView_, 1);

    QHBoxLayout *buttonLayout = new QHBoxLayout();
    buildButton_ = new QPushButton("Build");
    stopButton_ = new QPushButton("Stop");
    clearButton_ = new QPushButton("Clear");
    stopButton_->setEnabled(false);

    buttonLayout->addWidget(buildButton_);
    buttonLayout->addWidget(stopButton_);
    buttonLayout->addWidget(clearButton_);
    buttonLayout->addStretch();

    statusLabel_ = new QLabel("Ready");
    buttonLayout->addWidget(statusLabel_);

    mainLayout->addLayout(buttonLayout);

    connect(buildButton_, &QPushButton::clicked, this, &KernelBuildWidget::startBuild);
    connect(stopButton_, &QPushButton::clicked, this, &KernelBuildWidget::stopBuild);
    connect(clearButton_, &QPushButton::clicked, this, &KernelBuildWidget::clearOutput);
}

void KernelBuildWidget::connectSignals()
{
    if (!grpcClient_) {
        return;
    }

    connect(grpcClient_, &GrpcClient::buildStageChanged, this, &KernelBuildWidget::onBuildStageChanged);
    connect(grpcClient_, &GrpcClient::buildLogReceived, this, &KernelBuildWidget::onBuildLogReceived);
    connect(grpcClient_, &GrpcClient::buildErrorReceived, this, &KernelBuildWidget::onBuildErrorReceived);
    connect(grpcClient_, &GrpcClient::buildCompleted, this, &KernelBuildWidget::onBuildCompleted);
    connect(grpcClient_, &GrpcClient::errorOccurred, this, &KernelBuildWidget::onErrorOccurred);
}

void KernelBuildWidget::startBuild()
{
    if (isBuilding_ || !grpcClient_) {
        return;
    }

    isBuilding_ = true;
    buildButton_->setEnabled(false);
    stopButton_->setEnabled(true);
    progressBar_->setValue(0);
    statusLabel_->setText("Building...");
    outputView_->clear();

    QString arch = archCombo_->currentText();
    QString target = targetCombo_->currentText();
    int jobs = jobsSpinBox_->value();
    bool verbose = verboseCheckBox_->isChecked();

    QStringList targets;
    targets << target;

    QThread *buildThread = QThread::create([this, targets, jobs, arch, verbose]() {
        grpcClient_->buildKernel(targets, jobs, arch, verbose);
    });
    
    connect(buildThread, &QThread::finished, buildThread, &QThread::deleteLater);
    buildThread->start();
}

void KernelBuildWidget::stopBuild()
{
    isBuilding_ = false;
    buildButton_->setEnabled(true);
    stopButton_->setEnabled(false);
    statusLabel_->setText("Build stopped");
}

void KernelBuildWidget::clearOutput()
{
    outputView_->clear();
    progressBar_->setValue(0);
    stageLabel_->setText("Ready");
}

void KernelBuildWidget::onBuildStageChanged(const QString &stageName, int progress, int currentFile, int totalFiles)
{
    progressBar_->setValue(progress);
    
    QString stageText = stageName;
    if (totalFiles > 0) {
        stageText += QString(" (%1/%2)").arg(currentFile).arg(totalFiles);
    } else {
        stageText += QString(" (%1%)").arg(progress);
    }
    stageLabel_->setText(stageText);

    QString logText = QString("<span style='color: #4ec9b0;'>[STAGE]</span> %1 - %2%")
                         .arg(stageName)
                         .arg(progress);
    outputView_->append(logText);
    
    QScrollBar *scrollBar = outputView_->verticalScrollBar();
    scrollBar->setValue(scrollBar->maximum());
}

void KernelBuildWidget::onBuildLogReceived(int level, const QString &message, qint64 timestamp)
{
    QString levelStr = logLevelToString(level);
    QString color = logLevelToColor(level);
    QString timeStr = formatTimestamp(timestamp);
    
    QString logText = QString("<span style='color: #808080;'>[%1]</span> "
                             "<span style='color: %2;'>[%3]</span> %4")
                         .arg(timeStr)
                         .arg(color)
                         .arg(levelStr)
                         .arg(message.toHtmlEscaped());
    
    outputView_->append(logText);
    
    QScrollBar *scrollBar = outputView_->verticalScrollBar();
    scrollBar->setValue(scrollBar->maximum());
}

void KernelBuildWidget::onBuildErrorReceived(const QString &message, const QString &file, int line)
{
    QString errorText;
    if (!file.isEmpty() && line > 0) {
        errorText = QString("<span style='color: #f48771;'>[ERROR]</span> "
                           "<span style='color: #ce9178;'>%1:%2</span> %3")
                       .arg(file)
                       .arg(line)
                       .arg(message.toHtmlEscaped());
    } else {
        errorText = QString("<span style='color: #f48771;'>[ERROR]</span> %1")
                       .arg(message.toHtmlEscaped());
    }
    
    outputView_->append(errorText);
    
    QScrollBar *scrollBar = outputView_->verticalScrollBar();
    scrollBar->setValue(scrollBar->maximum());
}

void KernelBuildWidget::onBuildCompleted(bool success, qint64 durationMs, const QString &imagePath)
{
    isBuilding_ = false;
    buildButton_->setEnabled(true);
    stopButton_->setEnabled(false);
    progressBar_->setValue(success ? 100 : progressBar_->value());

    QString statusText;
    QString color;
    if (success) {
        statusText = QString("Build completed in %1s").arg(durationMs / 1000.0, 0, 'f', 2);
        statusLabel_->setText(statusText);
        color = "#4ec9b0";
        
        if (!imagePath.isEmpty()) {
            QString logText = QString("<span style='color: %1;'>[SUCCESS]</span> Image: %2")
                                 .arg(color)
                                 .arg(imagePath);
            outputView_->append(logText);
        }
    } else {
        statusText = QString("Build failed after %1s").arg(durationMs / 1000.0, 0, 'f', 2);
        statusLabel_->setText(statusText);
        color = "#f48771";
    }

    QString completionText = QString("<span style='color: %1;'>========== %2 ==========</span>")
                                .arg(color)
                                .arg(statusText);
    outputView_->append(completionText);
    
    QScrollBar *scrollBar = outputView_->verticalScrollBar();
    scrollBar->setValue(scrollBar->maximum());
}

void KernelBuildWidget::onErrorOccurred(const QString &error)
{
    isBuilding_ = false;
    buildButton_->setEnabled(true);
    stopButton_->setEnabled(false);
    statusLabel_->setText("Error occurred");

    QString errorText = QString("<span style='color: #f48771;'>[GRPC ERROR]</span> %1")
                           .arg(error.toHtmlEscaped());
    outputView_->append(errorText);
    
    QScrollBar *scrollBar = outputView_->verticalScrollBar();
    scrollBar->setValue(scrollBar->maximum());
}

QString KernelBuildWidget::formatTimestamp(qint64 timestampMs) const
{
    QDateTime dt = QDateTime::fromMSecsSinceEpoch(timestampMs);
    return dt.toString("HH:mm:ss");
}

QString KernelBuildWidget::logLevelToString(int level) const
{
    switch (level) {
        case 0: return "DEBUG";
        case 1: return "INFO";
        case 2: return "WARN";
        case 3: return "ERROR";
        default: return "UNKNOWN";
    }
}

QString KernelBuildWidget::logLevelToColor(int level) const
{
    switch (level) {
        case 0: return "#808080";
        case 1: return "#4ec9b0";
        case 2: return "#dcdcaa";
        case 3: return "#f48771";
        default: return "#d4d4d4";
    }
}
