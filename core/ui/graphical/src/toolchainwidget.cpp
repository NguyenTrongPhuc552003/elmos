#include "toolchainwidget.h"
#include "grpcclient.h"
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QGroupBox>
#include <QSplitter>
#include <QMessageBox>
#include <QThread>

ToolchainWidget::ToolchainWidget(GrpcClient *grpcClient, QWidget *parent)
    : QWidget(parent)
    , grpcClient_(grpcClient)
    , isBuilding_(false)
{
    setupUI();
    connectSignals();
    refreshStatus();
    refreshTargets();
}

ToolchainWidget::~ToolchainWidget() = default;

void ToolchainWidget::setupUI()
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);

    QSplitter *splitter = new QSplitter(Qt::Horizontal);

    QWidget *leftPanel = new QWidget();
    QVBoxLayout *leftLayout = new QVBoxLayout(leftPanel);

    QLabel *targetsLabel = new QLabel("Available Toolchain Targets:");
    leftLayout->addWidget(targetsLabel);

    targetsListWidget_ = new QListWidget();
    leftLayout->addWidget(targetsListWidget_, 1);

    connect(targetsListWidget_, &QListWidget::itemClicked, this, &ToolchainWidget::onTargetSelected);

    splitter->addWidget(leftPanel);

    QWidget *rightPanel = new QWidget();
    QVBoxLayout *rightLayout = new QVBoxLayout(rightPanel);

    QGroupBox *detailsGroup = new QGroupBox("Toolchain Details");
    QVBoxLayout *detailsLayout = new QVBoxLayout(detailsGroup);

    selectedTargetLabel_ = new QLabel("Target: Not selected");
    gccVersionLabel_ = new QLabel("GCC Version: N/A");
    pathLabel_ = new QLabel("Path: N/A");
    statusIconLabel_ = new QLabel("Status: ❓");

    detailsLayout->addWidget(selectedTargetLabel_);
    detailsLayout->addWidget(gccVersionLabel_);
    detailsLayout->addWidget(pathLabel_);
    detailsLayout->addWidget(statusIconLabel_);

    detailsTextEdit_ = new QTextEdit();
    detailsTextEdit_->setReadOnly(true);
    detailsTextEdit_->setMaximumHeight(150);
    detailsLayout->addWidget(detailsTextEdit_);

    rightLayout->addWidget(detailsGroup);

    statusLabel_ = new QLabel("Ready");
    rightLayout->addWidget(statusLabel_);

    progressBar_ = new QProgressBar();
    progressBar_->setRange(0, 100);
    progressBar_->setValue(0);
    progressBar_->setVisible(false);
    rightLayout->addWidget(progressBar_);

    QHBoxLayout *buttonLayout = new QHBoxLayout();
    installButton_ = new QPushButton("Install crosstool-ng");
    buildButton_ = new QPushButton("Build Toolchain");
    cleanButton_ = new QPushButton("Clean");
    refreshButton_ = new QPushButton("Refresh");

    buildButton_->setEnabled(false);
    cleanButton_->setEnabled(false);

    connect(installButton_, &QPushButton::clicked, this, &ToolchainWidget::installCrosstoolNG);
    connect(buildButton_, &QPushButton::clicked, this, &ToolchainWidget::buildSelectedToolchain);
    connect(cleanButton_, &QPushButton::clicked, this, &ToolchainWidget::cleanToolchain);
    connect(refreshButton_, &QPushButton::clicked, this, &ToolchainWidget::refreshStatus);

    buttonLayout->addWidget(installButton_);
    buttonLayout->addWidget(buildButton_);
    buttonLayout->addWidget(cleanButton_);
    buttonLayout->addWidget(refreshButton_);
    buttonLayout->addStretch();

    rightLayout->addLayout(buttonLayout);
    rightLayout->addStretch();

    splitter->addWidget(rightPanel);
    splitter->setStretchFactor(0, 1);
    splitter->setStretchFactor(1, 2);

    mainLayout->addWidget(splitter);
}

void ToolchainWidget::connectSignals()
{
    if (!grpcClient_) {
        return;
    }
}

void ToolchainWidget::refreshTargets()
{
    if (!grpcClient_) {
        return;
    }

    targetsListWidget_->clear();

    QStringList targets = {
        "aarch64-unknown-linux-gnu (ARM64)",
        "arm-cortex-a15-linux-gnueabihf (ARM 32-bit)",
        "riscv64-unknown-linux-gnu (RISC-V 64-bit)",
        "x86_64-unknown-linux-gnu (x86_64)"
    };

    targetsListWidget_->addItems(targets);
}

void ToolchainWidget::refreshStatus()
{
    if (!grpcClient_) {
        return;
    }

    statusLabel_->setText("Checking toolchain status...");
}

void ToolchainWidget::installCrosstoolNG()
{
    if (!grpcClient_) {
        QMessageBox::warning(this, "Error", "Not connected to ELMOS server");
        return;
    }

    isBuilding_ = true;
    installButton_->setEnabled(false);
    progressBar_->setVisible(true);
    progressBar_->setValue(0);
    statusLabel_->setText("Installing crosstool-ng...");

    QThread::create([this]() {
        detailsTextEdit_->append("Starting crosstool-ng installation...");
        detailsTextEdit_->append("This may take 5-10 minutes...");
    })->start();
}

void ToolchainWidget::buildSelectedToolchain()
{
    if (selectedTarget_.isEmpty()) {
        QMessageBox::warning(this, "No Target Selected", "Please select a toolchain target first.");
        return;
    }

    if (!grpcClient_) {
        QMessageBox::warning(this, "Error", "Not connected to ELMOS server");
        return;
    }

    QMessageBox::StandardButton reply = QMessageBox::question(
        this,
        "Build Toolchain",
        QString("Build toolchain for %1?\n\nThis will take 30-60 minutes.").arg(selectedTarget_),
        QMessageBox::Yes | QMessageBox::No
    );

    if (reply == QMessageBox::No) {
        return;
    }

    isBuilding_ = true;
    buildButton_->setEnabled(false);
    progressBar_->setVisible(true);
    progressBar_->setValue(0);
    statusLabel_->setText(QString("Building %1...").arg(selectedTarget_));
    detailsTextEdit_->append(QString("\n=== Building %1 ===").arg(selectedTarget_));

    QThread::create([this]() {
        detailsTextEdit_->append("Configuring crosstool-ng...");
        detailsTextEdit_->append("Building toolchain (this will take a while)...");
    })->start();
}

void ToolchainWidget::cleanToolchain()
{
    if (selectedTarget_.isEmpty()) {
        QMessageBox::warning(this, "No Target Selected", "Please select a toolchain target first.");
        return;
    }

    QMessageBox::StandardButton reply = QMessageBox::question(
        this,
        "Clean Toolchain",
        QString("Clean build artifacts for %1?").arg(selectedTarget_),
        QMessageBox::Yes | QMessageBox::No
    );

    if (reply == QMessageBox::Yes) {
        statusLabel_->setText(QString("Cleaning %1...").arg(selectedTarget_));
        detailsTextEdit_->append(QString("Cleaned %1").arg(selectedTarget_));
    }
}

void ToolchainWidget::onTargetSelected(QListWidgetItem *item)
{
    selectedTarget_ = item->text().split(" ").first();
    selectedTargetLabel_->setText(QString("Target: %1").arg(selectedTarget_));
    buildButton_->setEnabled(true);
    cleanButton_->setEnabled(true);

    detailsTextEdit_->clear();
    detailsTextEdit_->append(QString("Selected: %1").arg(selectedTarget_));
    detailsTextEdit_->append("");
    detailsTextEdit_->append("Configuration:");
    detailsTextEdit_->append(QString("  - Target: %1").arg(selectedTarget_));
    detailsTextEdit_->append("  - Toolchain: crosstool-ng");
    detailsTextEdit_->append("  - Features: GCC 13+, glibc, binutils");

    if (selectedTarget_.startsWith("aarch64")) {
        gccVersionLabel_->setText("GCC Version: 13.2.0 (estimated)");
        pathLabel_->setText("Path: ~/x-tools/aarch64-unknown-linux-gnu/");
    } else if (selectedTarget_.startsWith("arm")) {
        gccVersionLabel_->setText("GCC Version: 13.2.0 (estimated)");
        pathLabel_->setText("Path: ~/x-tools/arm-cortex-a15-linux-gnueabihf/");
    } else if (selectedTarget_.startsWith("riscv")) {
        gccVersionLabel_->setText("GCC Version: 13.2.0 (estimated)");
        pathLabel_->setText("Path: ~/x-tools/riscv64-unknown-linux-gnu/");
    }

    statusIconLabel_->setText("Status: ⏳ Not built");
}

void ToolchainWidget::onInstallProgress(int stage, int progress, const QString &message)
{
    progressBar_->setValue(progress);
    statusLabel_->setText(message);
    detailsTextEdit_->append(message);
}

void ToolchainWidget::onBuildProgress(int progress, const QString &message)
{
    progressBar_->setValue(progress);
    detailsTextEdit_->append(message);
}
