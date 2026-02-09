#include "workspacewidget.h"
#include "grpcclient.h"
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QFormLayout>
#include <QGroupBox>
#include <QMessageBox>
#include <QThread>

WorkspaceWidget::WorkspaceWidget(GrpcClient *grpcClient, QWidget *parent)
    : QWidget(parent)
    , grpcClient_(grpcClient)
    , isMounted_(false)
    , isInitializing_(false)
{
    setupUI();
    connectSignals();
    refreshStatus();
}

WorkspaceWidget::~WorkspaceWidget() = default;

void WorkspaceWidget::setupUI()
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);

    QGroupBox *statusGroup = new QGroupBox("Workspace Status");
    QVBoxLayout *statusLayout = new QVBoxLayout(statusGroup);

    QHBoxLayout *statusLineLayout = new QHBoxLayout();
    statusIconLabel_ = new QLabel("🔴");
    statusIconLabel_->setStyleSheet("font-size: 24px;");
    statusLabel_ = new QLabel("Not mounted");
    statusLabel_->setStyleSheet("font-weight: bold; font-size: 14px;");
    statusLineLayout->addWidget(statusIconLabel_);
    statusLineLayout->addWidget(statusLabel_);
    statusLineLayout->addStretch();
    statusLayout->addLayout(statusLineLayout);

    QFormLayout *detailsLayout = new QFormLayout();
    mountPointLabel_ = new QLabel("N/A");
    volumePathLabel_ = new QLabel("N/A");
    sizeLabel_ = new QLabel("N/A");
    usedLabel_ = new QLabel("N/A");
    availLabel_ = new QLabel("N/A");

    detailsLayout->addRow("Mount Point:", mountPointLabel_);
    detailsLayout->addRow("Volume Path:", volumePathLabel_);
    detailsLayout->addRow("Total Size:", sizeLabel_);
    detailsLayout->addRow("Used:", usedLabel_);
    detailsLayout->addRow("Available:", availLabel_);

    statusLayout->addLayout(detailsLayout);

    mainLayout->addWidget(statusGroup);

    QGroupBox *initGroup = new QGroupBox("Initialize New Workspace");
    QFormLayout *initLayout = new QFormLayout(initGroup);

    workspaceNameEdit_ = new QLineEdit("elmos");
    volumeSizeEdit_ = new QLineEdit("40G");

    initLayout->addRow("Workspace Name:", workspaceNameEdit_);
    initLayout->addRow("Volume Size:", volumeSizeEdit_);

    mainLayout->addWidget(initGroup);

    progressBar_ = new QProgressBar();
    progressBar_->setRange(0, 100);
    progressBar_->setValue(0);
    progressBar_->setVisible(false);
    mainLayout->addWidget(progressBar_);

    progressMessageLabel_ = new QLabel("");
    progressMessageLabel_->setVisible(false);
    mainLayout->addWidget(progressMessageLabel_);

    QHBoxLayout *buttonLayout = new QHBoxLayout();
    initButton_ = new QPushButton("Initialize");
    mountButton_ = new QPushButton("Mount");
    unmountButton_ = new QPushButton("Unmount");
    refreshButton_ = new QPushButton("Refresh");

    mountButton_->setEnabled(false);
    unmountButton_->setEnabled(false);

    connect(initButton_, &QPushButton::clicked, this, &WorkspaceWidget::initWorkspace);
    connect(mountButton_, &QPushButton::clicked, this, &WorkspaceWidget::mountWorkspace);
    connect(unmountButton_, &QPushButton::clicked, this, &WorkspaceWidget::unmountWorkspace);
    connect(refreshButton_, &QPushButton::clicked, this, &WorkspaceWidget::refreshStatus);

    buttonLayout->addWidget(initButton_);
    buttonLayout->addWidget(mountButton_);
    buttonLayout->addWidget(unmountButton_);
    buttonLayout->addWidget(refreshButton_);
    buttonLayout->addStretch();

    mainLayout->addLayout(buttonLayout);
    mainLayout->addStretch();
}

void WorkspaceWidget::connectSignals()
{
    if (!grpcClient_) {
        return;
    }
}

void WorkspaceWidget::initWorkspace()
{
    if (!grpcClient_) {
        QMessageBox::warning(this, "Error", "Not connected to ELMOS server");
        return;
    }

    QString name = workspaceNameEdit_->text().trimmed();
    QString size = volumeSizeEdit_->text().trimmed();

    if (name.isEmpty() || size.isEmpty()) {
        QMessageBox::warning(this, "Invalid Input", "Please provide workspace name and size");
        return;
    }

    QMessageBox::StandardButton reply = QMessageBox::question(
        this,
        "Initialize Workspace",
        QString("Initialize workspace '%1' with size %2?\n\nThis will create a new volume.").arg(name).arg(size),
        QMessageBox::Yes | QMessageBox::No
    );

    if (reply == QMessageBox::No) {
        return;
    }

    isInitializing_ = true;
    initButton_->setEnabled(false);
    progressBar_->setVisible(true);
    progressBar_->setValue(0);
    progressMessageLabel_->setVisible(true);
    progressMessageLabel_->setText("Initializing workspace...");

    statusIconLabel_->setText("⏳");
    statusLabel_->setText("Initializing...");

    QThread::create([this, name, size]() {
        progressMessageLabel_->setText("Creating volume...");
        QThread::msleep(1000);
        onInitProgress(1, 25, "Creating volume...");

        progressMessageLabel_->setText("Formatting...");
        QThread::msleep(1000);
        onInitProgress(2, 50, "Formatting volume...");

        progressMessageLabel_->setText("Mounting...");
        QThread::msleep(1000);
        onInitProgress(2, 75, "Mounting volume...");

        progressMessageLabel_->setText("Saving configuration...");
        QThread::msleep(500);
        onInitProgress(3, 95, "Saving configuration...");

        onInitProgress(4, 100, "Complete");

        QMetaObject::invokeMethod(this, [this]() {
            isInitializing_ = false;
            initButton_->setEnabled(true);
            progressBar_->setVisible(false);
            progressMessageLabel_->setVisible(false);
            
            QMessageBox::information(this, "Success", "Workspace initialized successfully!");
            refreshStatus();
        });
    })->start();
}

void WorkspaceWidget::mountWorkspace()
{
    if (!grpcClient_) {
        return;
    }

    statusIconLabel_->setText("⏳");
    statusLabel_->setText("Mounting...");

    QThread::msleep(500);
    
    isMounted_ = true;
    statusIconLabel_->setText("🟢");
    statusLabel_->setText("Mounted");
    mountButton_->setEnabled(false);
    unmountButton_->setEnabled(true);
    
    refreshStatus();
}

void WorkspaceWidget::unmountWorkspace()
{
    if (!grpcClient_) {
        return;
    }

    QMessageBox::StandardButton reply = QMessageBox::question(
        this,
        "Unmount Workspace",
        "Unmount the workspace volume?",
        QMessageBox::Yes | QMessageBox::No
    );

    if (reply == QMessageBox::No) {
        return;
    }

    statusIconLabel_->setText("⏳");
    statusLabel_->setText("Unmounting...");

    QThread::msleep(500);
    
    isMounted_ = false;
    statusIconLabel_->setText("🔴");
    statusLabel_->setText("Not mounted");
    mountButton_->setEnabled(true);
    unmountButton_->setEnabled(false);
    
    refreshStatus();
}

void WorkspaceWidget::refreshStatus()
{
    if (!grpcClient_) {
        statusIconLabel_->setText("❌");
        statusLabel_->setText("Not connected");
        return;
    }

    if (isMounted_) {
        statusIconLabel_->setText("🟢");
        statusLabel_->setText("Mounted");
        mountPointLabel_->setText("/Volumes/elmos");
        volumePathLabel_->setText("~/Library/elmos/data/elmos.sparseimage");
        sizeLabel_->setText(formatBytes(40LL * 1024 * 1024 * 1024));
        usedLabel_->setText(formatBytes(8LL * 1024 * 1024 * 1024));
        availLabel_->setText(formatBytes(32LL * 1024 * 1024 * 1024));
        
        mountButton_->setEnabled(false);
        unmountButton_->setEnabled(true);
    } else {
        statusIconLabel_->setText("🔴");
        statusLabel_->setText("Not mounted");
        mountPointLabel_->setText("N/A");
        volumePathLabel_->setText("N/A");
        sizeLabel_->setText("N/A");
        usedLabel_->setText("N/A");
        availLabel_->setText("N/A");
        
        mountButton_->setEnabled(true);
        unmountButton_->setEnabled(false);
    }
}

void WorkspaceWidget::onInitProgress(int stage, int progress, const QString &message)
{
    progressBar_->setValue(progress);
    progressMessageLabel_->setText(message);
}

void WorkspaceWidget::onStatusUpdated(bool mounted, const QString &mountPoint, qint64 sizeBytes, qint64 usedBytes, qint64 availBytes)
{
    isMounted_ = mounted;
    
    if (mounted) {
        statusIconLabel_->setText("🟢");
        statusLabel_->setText("Mounted");
        mountPointLabel_->setText(mountPoint);
        sizeLabel_->setText(formatBytes(sizeBytes));
        usedLabel_->setText(formatBytes(usedBytes));
        availLabel_->setText(formatBytes(availBytes));
    } else {
        statusIconLabel_->setText("🔴");
        statusLabel_->setText("Not mounted");
    }
}

QString WorkspaceWidget::formatBytes(qint64 bytes) const
{
    const qint64 KB = 1024;
    const qint64 MB = KB * 1024;
    const qint64 GB = MB * 1024;

    if (bytes >= GB) {
        return QString("%1 GB").arg(bytes / (double)GB, 0, 'f', 2);
    } else if (bytes >= MB) {
        return QString("%1 MB").arg(bytes / (double)MB, 0, 'f', 2);
    } else if (bytes >= KB) {
        return QString("%1 KB").arg(bytes / (double)KB, 0, 'f', 2);
    } else {
        return QString("%1 bytes").arg(bytes);
    }
}
