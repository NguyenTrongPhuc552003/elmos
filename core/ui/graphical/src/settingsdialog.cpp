#include "settingsdialog.h"
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QFormLayout>
#include <QGroupBox>
#include <QLabel>
#include <QDialogButtonBox>
#include <QMessageBox>

SettingsDialog::SettingsDialog(QWidget *parent)
    : QDialog(parent)
{
    setWindowTitle("ELMOS Settings");
    resize(600, 500);
    setupUI();
    loadSettings();
}

SettingsDialog::~SettingsDialog() = default;

void SettingsDialog::setupUI()
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);

    tabWidget_ = new QTabWidget();

    QWidget *editorTab = new QWidget();
    QFormLayout *editorLayout = new QFormLayout(editorTab);

    editorFontCombo_ = new QFontComboBox();
    editorFontSizeSpinBox_ = new QSpinBox();
    editorFontSizeSpinBox_->setRange(8, 24);
    editorFontSizeSpinBox_->setValue(11);

    editorLineNumbersCheckBox_ = new QCheckBox("Show line numbers");
    editorLineNumbersCheckBox_->setChecked(true);
    editorHighlightCurrentLineCheckBox_ = new QCheckBox("Highlight current line");
    editorHighlightCurrentLineCheckBox_->setChecked(true);
    editorAutoIndentCheckBox_ = new QCheckBox("Auto-indent");
    editorAutoIndentCheckBox_->setChecked(true);

    editorTabWidthSpinBox_ = new QSpinBox();
    editorTabWidthSpinBox_->setRange(2, 8);
    editorTabWidthSpinBox_->setValue(4);

    editorLayout->addRow("Font:", editorFontCombo_);
    editorLayout->addRow("Font Size:", editorFontSizeSpinBox_);
    editorLayout->addRow("Tab Width:", editorTabWidthSpinBox_);
    editorLayout->addRow(editorLineNumbersCheckBox_);
    editorLayout->addRow(editorHighlightCurrentLineCheckBox_);
    editorLayout->addRow(editorAutoIndentCheckBox_);

    tabWidget_->addTab(editorTab, "Editor");

    QWidget *grpcTab = new QWidget();
    QFormLayout *grpcLayout = new QFormLayout(grpcTab);

    grpcServerEdit_ = new QLineEdit("unix:///tmp/elmos.sock");
    grpcTimeoutSpinBox_ = new QSpinBox();
    grpcTimeoutSpinBox_->setRange(5, 300);
    grpcTimeoutSpinBox_->setValue(30);
    grpcTimeoutSpinBox_->setSuffix(" seconds");
    grpcAutoConnectCheckBox_ = new QCheckBox("Auto-connect on startup");
    grpcAutoConnectCheckBox_->setChecked(true);

    grpcLayout->addRow("Server Address:", grpcServerEdit_);
    grpcLayout->addRow("Timeout:", grpcTimeoutSpinBox_);
    grpcLayout->addRow(grpcAutoConnectCheckBox_);

    tabWidget_->addTab(grpcTab, "Connection");

    QWidget *buildTab = new QWidget();
    QFormLayout *buildLayout = new QFormLayout(buildTab);

    buildArchCombo_ = new QComboBox();
    buildArchCombo_->addItems({"arm64", "arm", "riscv", "x86_64"});

    buildJobsSpinBox_ = new QSpinBox();
    buildJobsSpinBox_->setRange(0, 128);
    buildJobsSpinBox_->setValue(0);
    buildJobsSpinBox_->setSpecialValueText("Auto");

    buildVerboseCheckBox_ = new QCheckBox("Verbose build output");
    buildUseLLVMCheckBox_ = new QCheckBox("Use LLVM toolchain");
    buildUseLLVMCheckBox_->setChecked(true);

    buildToolchainPathEdit_ = new QLineEdit();
    buildToolchainPathEdit_->setPlaceholderText("/opt/homebrew/opt/llvm");

    buildLayout->addRow("Default Architecture:", buildArchCombo_);
    buildLayout->addRow("Default Jobs:", buildJobsSpinBox_);
    buildLayout->addRow("Toolchain Path:", buildToolchainPathEdit_);
    buildLayout->addRow(buildVerboseCheckBox_);
    buildLayout->addRow(buildUseLLVMCheckBox_);

    tabWidget_->addTab(buildTab, "Build");

    QWidget *qemuTab = new QWidget();
    QFormLayout *qemuLayout = new QFormLayout(qemuTab);

    qemuMemorySpinBox_ = new QSpinBox();
    qemuMemorySpinBox_->setRange(64, 16384);
    qemuMemorySpinBox_->setValue(256);
    qemuMemorySpinBox_->setSuffix(" MB");

    qemuCPUsSpinBox_ = new QSpinBox();
    qemuCPUsSpinBox_->setRange(1, 16);
    qemuCPUsSpinBox_->setValue(2);

    qemuMachineCombo_ = new QComboBox();
    qemuMachineCombo_->addItems({"virt", "raspi3", "raspi4", "pc", "q35"});

    qemuExtraArgsEdit_ = new QLineEdit();
    qemuExtraArgsEdit_->setPlaceholderText("-enable-kvm");

    qemuLayout->addRow("Memory:", qemuMemorySpinBox_);
    qemuLayout->addRow("CPUs:", qemuCPUsSpinBox_);
    qemuLayout->addRow("Machine:", qemuMachineCombo_);
    qemuLayout->addRow("Extra Arguments:", qemuExtraArgsEdit_);

    tabWidget_->addTab(qemuTab, "QEMU");

    QWidget *appearanceTab = new QWidget();
    QFormLayout *appearanceLayout = new QFormLayout(appearanceTab);

    themeCombo_ = new QComboBox();
    themeCombo_->addItems({"Dark (Default)", "Light", "System"});

    showStatusBarCheckBox_ = new QCheckBox("Show status bar");
    showStatusBarCheckBox_->setChecked(true);
    showToolBarCheckBox_ = new QCheckBox("Show toolbar");
    showToolBarCheckBox_->setChecked(true);
    restoreSessionCheckBox_ = new QCheckBox("Restore previous session");
    restoreSessionCheckBox_->setChecked(true);

    appearanceLayout->addRow("Theme:", themeCombo_);
    appearanceLayout->addRow(showStatusBarCheckBox_);
    appearanceLayout->addRow(showToolBarCheckBox_);
    appearanceLayout->addRow(restoreSessionCheckBox_);

    tabWidget_->addTab(appearanceTab, "Appearance");

    mainLayout->addWidget(tabWidget_);

    QHBoxLayout *buttonLayout = new QHBoxLayout();
    restoreDefaultsButton_ = new QPushButton("Restore Defaults");
    buttonLayout->addWidget(restoreDefaultsButton_);
    buttonLayout->addStretch();

    okButton_ = new QPushButton("OK");
    cancelButton_ = new QPushButton("Cancel");
    applyButton_ = new QPushButton("Apply");

    buttonLayout->addWidget(okButton_);
    buttonLayout->addWidget(cancelButton_);
    buttonLayout->addWidget(applyButton_);

    mainLayout->addLayout(buttonLayout);

    connect(okButton_, &QPushButton::clicked, this, &SettingsDialog::accept);
    connect(cancelButton_, &QPushButton::clicked, this, &SettingsDialog::reject);
    connect(applyButton_, &QPushButton::clicked, this, &SettingsDialog::saveSettings);
    connect(restoreDefaultsButton_, &QPushButton::clicked, this, &SettingsDialog::restoreDefaults);
}

void SettingsDialog::loadSettings()
{
    QSettings settings("ELMOS", "ELMOS-IDE");

    editorFontCombo_->setCurrentFont(QFont(settings.value("editor/font", "Monaco").toString()));
    editorFontSizeSpinBox_->setValue(settings.value("editor/fontSize", 11).toInt());
    editorTabWidthSpinBox_->setValue(settings.value("editor/tabWidth", 4).toInt());
    editorLineNumbersCheckBox_->setChecked(settings.value("editor/showLineNumbers", true).toBool());
    editorHighlightCurrentLineCheckBox_->setChecked(settings.value("editor/highlightCurrentLine", true).toBool());
    editorAutoIndentCheckBox_->setChecked(settings.value("editor/autoIndent", true).toBool());

    grpcServerEdit_->setText(settings.value("grpc/serverAddress", "unix:///tmp/elmos.sock").toString());
    grpcTimeoutSpinBox_->setValue(settings.value("grpc/timeout", 30).toInt());
    grpcAutoConnectCheckBox_->setChecked(settings.value("grpc/autoConnect", true).toBool());

    buildArchCombo_->setCurrentText(settings.value("build/arch", "arm64").toString());
    buildJobsSpinBox_->setValue(settings.value("build/jobs", 0).toInt());
    buildToolchainPathEdit_->setText(settings.value("build/toolchainPath", "/opt/homebrew/opt/llvm").toString());
    buildVerboseCheckBox_->setChecked(settings.value("build/verbose", false).toBool());
    buildUseLLVMCheckBox_->setChecked(settings.value("build/useLLVM", true).toBool());

    qemuMemorySpinBox_->setValue(settings.value("qemu/memory", 256).toInt());
    qemuCPUsSpinBox_->setValue(settings.value("qemu/cpus", 2).toInt());
    qemuMachineCombo_->setCurrentText(settings.value("qemu/machine", "virt").toString());
    qemuExtraArgsEdit_->setText(settings.value("qemu/extraArgs", "").toString());

    themeCombo_->setCurrentIndex(settings.value("appearance/theme", 0).toInt());
    showStatusBarCheckBox_->setChecked(settings.value("appearance/showStatusBar", true).toBool());
    showToolBarCheckBox_->setChecked(settings.value("appearance/showToolBar", true).toBool());
    restoreSessionCheckBox_->setChecked(settings.value("appearance/restoreSession", true).toBool());
}

void SettingsDialog::saveSettings()
{
    QSettings settings("ELMOS", "ELMOS-IDE");

    settings.setValue("editor/font", editorFontCombo_->currentFont().family());
    settings.setValue("editor/fontSize", editorFontSizeSpinBox_->value());
    settings.setValue("editor/tabWidth", editorTabWidthSpinBox_->value());
    settings.setValue("editor/showLineNumbers", editorLineNumbersCheckBox_->isChecked());
    settings.setValue("editor/highlightCurrentLine", editorHighlightCurrentLineCheckBox_->isChecked());
    settings.setValue("editor/autoIndent", editorAutoIndentCheckBox_->isChecked());

    settings.setValue("grpc/serverAddress", grpcServerEdit_->text());
    settings.setValue("grpc/timeout", grpcTimeoutSpinBox_->value());
    settings.setValue("grpc/autoConnect", grpcAutoConnectCheckBox_->isChecked());

    settings.setValue("build/arch", buildArchCombo_->currentText());
    settings.setValue("build/jobs", buildJobsSpinBox_->value());
    settings.setValue("build/toolchainPath", buildToolchainPathEdit_->text());
    settings.setValue("build/verbose", buildVerboseCheckBox_->isChecked());
    settings.setValue("build/useLLVM", buildUseLLVMCheckBox_->isChecked());

    settings.setValue("qemu/memory", qemuMemorySpinBox_->value());
    settings.setValue("qemu/cpus", qemuCPUsSpinBox_->value());
    settings.setValue("qemu/machine", qemuMachineCombo_->currentText());
    settings.setValue("qemu/extraArgs", qemuExtraArgsEdit_->text());

    settings.setValue("appearance/theme", themeCombo_->currentIndex());
    settings.setValue("appearance/showStatusBar", showStatusBarCheckBox_->isChecked());
    settings.setValue("appearance/showToolBar", showToolBarCheckBox_->isChecked());
    settings.setValue("appearance/restoreSession", restoreSessionCheckBox_->isChecked());

    settings.sync();
}

void SettingsDialog::accept()
{
    saveSettings();
    QDialog::accept();
}

void SettingsDialog::reject()
{
    QDialog::reject();
}

void SettingsDialog::restoreDefaults()
{
    QMessageBox::StandardButton reply = QMessageBox::question(
        this,
        "Restore Defaults",
        "Are you sure you want to restore all settings to defaults?",
        QMessageBox::Yes | QMessageBox::No
    );

    if (reply == QMessageBox::Yes) {
        QSettings settings("ELMOS", "ELMOS-IDE");
        settings.clear();
        loadSettings();
    }
}
