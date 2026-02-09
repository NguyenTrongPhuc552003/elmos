#include "archselectorwidget.h"
#include "grpcclient.h"
#include <QMessageBox>
#include <QProgressDialog>
#include <QDir>

ArchSelectorWidget::ArchSelectorWidget(GrpcClient *client, QWidget *parent)
    : QWidget(parent)
    , grpcClient_(client)
    , currentArch_("arm")
{
    archDescriptions_["arm"] = "ARM 32-bit (Cortex-A)";
    archDescriptions_["arm64"] = "ARM 64-bit (AArch64)";
    archDescriptions_["x86"] = "Intel x86 32-bit";
    archDescriptions_["x86_64"] = "Intel x86 64-bit (AMD64)";
    archDescriptions_["mips"] = "MIPS 32-bit";
    archDescriptions_["mips64"] = "MIPS 64-bit";
    archDescriptions_["powerpc"] = "PowerPC 32-bit";
    archDescriptions_["riscv64"] = "RISC-V 64-bit";
    
    setupUI();
    populateArchitectures();
    refreshToolchainStatus();
}

QString ArchSelectorWidget::currentArchitecture() const
{
    return currentArch_;
}

void ArchSelectorWidget::setArchitecture(const QString &arch)
{
    if (currentArch_ != arch) {
        int index = archComboBox_->findText(arch, Qt::MatchStartsWith);
        if (index >= 0) {
            archComboBox_->setCurrentIndex(index);
        }
    }
}

void ArchSelectorWidget::setupUI()
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);
    
    QGroupBox *selectorGroup = new QGroupBox(tr("Target Architecture"), this);
    QVBoxLayout *groupLayout = new QVBoxLayout(selectorGroup);
    
    QLabel *label = new QLabel(tr("Select target architecture for cross-compilation:"), this);
    groupLayout->addWidget(label);
    
    archComboBox_ = new QComboBox(this);
    archComboBox_->setSizeAdjustPolicy(QComboBox::AdjustToContents);
    connect(archComboBox_, QOverload<int>::of(&QComboBox::currentIndexChanged),
            this, &ArchSelectorWidget::onArchitectureSelected);
    groupLayout->addWidget(archComboBox_);
    
    statusLabel_ = new QLabel(this);
    statusLabel_->setWordWrap(true);
    statusLabel_->setStyleSheet("QLabel { padding: 8px; background: #f0f0f0; border-radius: 4px; }");
    groupLayout->addWidget(statusLabel_);
    
    detailsLabel_ = new QLabel(this);
    detailsLabel_->setWordWrap(true);
    detailsLabel_->setStyleSheet("QLabel { padding: 4px; color: #666; font-size: 10pt; }");
    groupLayout->addWidget(detailsLabel_);
    
    QHBoxLayout *buttonLayout = new QHBoxLayout();
    
    installButton_ = new QPushButton(tr("Install Toolchain"), this);
    installButton_->setIcon(QIcon::fromTheme("document-save"));
    connect(installButton_, &QPushButton::clicked, this, &ArchSelectorWidget::installToolchainForCurrentArch);
    buttonLayout->addWidget(installButton_);
    
    refreshButton_ = new QPushButton(tr("Refresh Status"), this);
    refreshButton_->setIcon(QIcon::fromTheme("view-refresh"));
    connect(refreshButton_, &QPushButton::clicked, this, &ArchSelectorWidget::refreshToolchainStatus);
    buttonLayout->addWidget(refreshButton_);
    
    detailsButton_ = new QPushButton(tr("Architecture Details"), this);
    detailsButton_->setIcon(QIcon::fromTheme("help-about"));
    connect(detailsButton_, &QPushButton::clicked, this, &ArchSelectorWidget::showArchitectureDetails);
    buttonLayout->addWidget(detailsButton_);
    
    buttonLayout->addStretch();
    groupLayout->addLayout(buttonLayout);
    
    mainLayout->addWidget(selectorGroup);
    mainLayout->addStretch();
}

void ArchSelectorWidget::populateArchitectures()
{
    archComboBox_->clear();
    
    QStringList architectures;
    architectures << "arm" << "arm64" << "x86" << "x86_64" 
                  << "mips" << "mips64" << "powerpc" << "riscv64";
    
    for (const QString &arch : architectures) {
        QString displayText = arch;
        if (archDescriptions_.contains(arch)) {
            displayText += " - " + archDescriptions_[arch];
        }
        archComboBox_->addItem(displayText, arch);
    }
    
    int armIndex = archComboBox_->findData("arm");
    if (armIndex >= 0) {
        archComboBox_->setCurrentIndex(armIndex);
    }
}

void ArchSelectorWidget::onArchitectureSelected(int index)
{
    if (index < 0) return;
    
    QString newArch = archComboBox_->itemData(index).toString();
    if (newArch.isEmpty()) {
        newArch = archComboBox_->itemText(index).split(" - ").first();
    }
    
    if (newArch != currentArch_) {
        currentArch_ = newArch;
        updateToolchainStatus(currentArch_);
        updateStatusDisplay();
        emit architectureChanged(currentArch_);
    }
}

void ArchSelectorWidget::refreshToolchainStatus()
{
    toolchainStatus_.clear();
    
    QStringList architectures;
    architectures << "arm" << "arm64" << "x86" << "x86_64" 
                  << "mips" << "mips64" << "powerpc" << "riscv64";
    
    for (const QString &arch : architectures) {
        updateToolchainStatus(arch);
    }
    
    updateStatusDisplay();
}

void ArchSelectorWidget::updateToolchainStatus(const QString &arch)
{
    bool installed = false;
    
    QString toolchainPath = QString("/opt/elmos/toolchains/%1").arg(arch);
    QDir toolchainDir(toolchainPath);
    installed = toolchainDir.exists();
    
    toolchainStatus_[arch] = installed;
    emit toolchainStatusUpdated(arch, installed);
}

void ArchSelectorWidget::updateStatusDisplay()
{
    bool installed = toolchainStatus_.value(currentArch_, false);
    
    if (installed) {
        statusLabel_->setText(QString("<b style='color:green;'>✓ Toolchain Installed</b><br>"
                                     "Target: %1").arg(currentArch_));
        statusLabel_->setStyleSheet("QLabel { padding: 8px; background: #d4edda; "
                                   "border: 1px solid #c3e6cb; border-radius: 4px; }");
        installButton_->setEnabled(false);
        installButton_->setText(tr("Toolchain Installed"));
    } else {
        statusLabel_->setText(QString("<b style='color:orange;'>⚠ Toolchain Not Installed</b><br>"
                                     "Target: %1<br>"
                                     "<small>Click 'Install Toolchain' to build it</small>").arg(currentArch_));
        statusLabel_->setStyleSheet("QLabel { padding: 8px; background: #fff3cd; "
                                   "border: 1px solid #ffc107; border-radius: 4px; }");
        installButton_->setEnabled(true);
        installButton_->setText(tr("Install Toolchain"));
    }
    
    if (archDescriptions_.contains(currentArch_)) {
        detailsLabel_->setText(tr("Architecture: %1").arg(archDescriptions_[currentArch_]));
    } else {
        detailsLabel_->setText(tr("Architecture: %1").arg(currentArch_));
    }
}

void ArchSelectorWidget::installToolchainForCurrentArch()
{
    QMessageBox::StandardButton reply = QMessageBox::question(
        this,
        tr("Install Toolchain"),
        tr("Building the %1 toolchain will take 20-40 minutes and requires:\n\n"
           "• 2-4 GB disk space\n"
           "• Active internet connection\n"
           "• Build tools (make, gcc)\n\n"
           "Do you want to proceed?").arg(currentArch_),
        QMessageBox::Yes | QMessageBox::No
    );
    
    if (reply != QMessageBox::Yes) {
        return;
    }
    
    QProgressDialog *progress = new QProgressDialog(
        tr("Building toolchain for %1...\n\nThis may take 20-40 minutes.").arg(currentArch_),
        tr("Cancel"),
        0, 0,
        this
    );
    progress->setWindowModality(Qt::WindowModal);
    progress->setMinimumDuration(0);
    progress->setValue(0);
    progress->show();
    
    QMessageBox::information(this, tr("Toolchain Build"),
                            tr("In a production implementation, this would:\n\n"
                               "1. Call gRPC BuildToolchain(%1)\n"
                               "2. Stream build progress\n"
                               "3. Update status on completion\n\n"
                               "For now, simulating installation...").arg(currentArch_));
    
    progress->close();
    delete progress;
    
    toolchainStatus_[currentArch_] = true;
    updateStatusDisplay();
    emit toolchainStatusUpdated(currentArch_, true);
}

void ArchSelectorWidget::showArchitectureDetails()
{
    QString details = tr("<h3>Architecture: %1</h3>").arg(currentArch_);
    
    if (archDescriptions_.contains(currentArch_)) {
        details += tr("<p><b>Description:</b> %1</p>").arg(archDescriptions_[currentArch_]);
    }
    
    details += tr("<p><b>Toolchain Status:</b> %1</p>")
                .arg(toolchainStatus_.value(currentArch_, false) ? 
                     tr("Installed ✓") : tr("Not Installed"));
    
    details += tr("<hr><p><b>Typical Use Cases:</b></p><ul>");
    
    if (currentArch_ == "arm") {
        details += tr("<li>Raspberry Pi (32-bit)</li>"
                     "<li>BeagleBone Black</li>"
                     "<li>32-bit ARM embedded systems</li>");
    } else if (currentArch_ == "arm64") {
        details += tr("<li>Raspberry Pi 3/4/5 (64-bit)</li>"
                     "<li>NVIDIA Jetson</li>"
                     "<li>Modern ARM servers</li>");
    } else if (currentArch_ == "x86_64") {
        details += tr("<li>Standard PC/Server</li>"
                     "<li>Virtual machines</li>"
                     "<li>Intel/AMD 64-bit systems</li>");
    } else if (currentArch_ == "riscv64") {
        details += tr("<li>RISC-V development boards</li>"
                     "<li>Open-source hardware platforms</li>");
    } else {
        details += tr("<li>Various embedded systems</li>"
                     "<li>Legacy hardware support</li>");
    }
    
    details += "</ul>";
    
    QMessageBox msgBox(this);
    msgBox.setWindowTitle(tr("Architecture Details"));
    msgBox.setTextFormat(Qt::RichText);
    msgBox.setText(details);
    msgBox.setIcon(QMessageBox::Information);
    msgBox.exec();
}
