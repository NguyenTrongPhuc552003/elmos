#ifndef SETTINGSDIALOG_H
#define SETTINGSDIALOG_H

#include <QDialog>
#include <QTabWidget>
#include <QLineEdit>
#include <QSpinBox>
#include <QCheckBox>
#include <QComboBox>
#include <QFontComboBox>
#include <QPushButton>
#include <QSettings>

class SettingsDialog : public QDialog
{
    Q_OBJECT

public:
    explicit SettingsDialog(QWidget *parent = nullptr);
    ~SettingsDialog();

public slots:
    void accept() override;
    void reject() override;
    void restoreDefaults();

private:
    void setupUI();
    void loadSettings();
    void saveSettings();

    QTabWidget *tabWidget_;

    QFontComboBox *editorFontCombo_;
    QSpinBox *editorFontSizeSpinBox_;
    QCheckBox *editorLineNumbersCheckBox_;
    QCheckBox *editorHighlightCurrentLineCheckBox_;
    QCheckBox *editorAutoIndentCheckBox_;
    QSpinBox *editorTabWidthSpinBox_;

    QLineEdit *grpcServerEdit_;
    QSpinBox *grpcTimeoutSpinBox_;
    QCheckBox *grpcAutoConnectCheckBox_;

    QComboBox *buildArchCombo_;
    QSpinBox *buildJobsSpinBox_;
    QCheckBox *buildVerboseCheckBox_;
    QCheckBox *buildUseLLVMCheckBox_;
    QLineEdit *buildToolchainPathEdit_;

    QSpinBox *qemuMemorySpinBox_;
    QSpinBox *qemuCPUsSpinBox_;
    QComboBox *qemuMachineCombo_;
    QLineEdit *qemuExtraArgsEdit_;

    QComboBox *themeCombo_;
    QCheckBox *showStatusBarCheckBox_;
    QCheckBox *showToolBarCheckBox_;
    QCheckBox *restoreSessionCheckBox_;

    QPushButton *restoreDefaultsButton_;
    QPushButton *okButton_;
    QPushButton *cancelButton_;
    QPushButton *applyButton_;
};

#endif
