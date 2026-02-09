#ifndef QEMUCONSOLEWIDGET_H
#define QEMUCONSOLEWIDGET_H

#include <QWidget>
#include <QTextEdit>
#include <QPushButton>
#include <QLabel>
#include <QSpinBox>
#include <QCheckBox>
#include <QLineEdit>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QTextCharFormat>

class GrpcClient;

class QEMUConsoleWidget : public QWidget
{
    Q_OBJECT

public:
    explicit QEMUConsoleWidget(GrpcClient *grpcClient, QWidget *parent = nullptr);
    ~QEMUConsoleWidget();

public slots:
    void startQEMU();
    void stopQEMU();
    void clearConsole();
    void sendInput();

private slots:
    void onQEMUStarted(int pid, const QString &qemuVersion, const QString &command);
    void onQEMUConsoleOutput(const QByteArray &data);
    void onQEMUStopped(int exitCode, qint64 uptimeMs);
    void onQEMUError(const QString &error);
    void onReturnPressed();

private:
    void setupUI();
    void connectSignals();
    void processANSI(const QByteArray &data);
    void appendText(const QString &text, const QTextCharFormat &format);

    GrpcClient *grpcClient_;
    bool isRunning_;
    QString inputBuffer_;

    QTextEdit *consoleView_;
    QLineEdit *inputLine_;
    QPushButton *startButton_;
    QPushButton *stopButton_;
    QPushButton *clearButton_;
    QLabel *statusLabel_;
    
    QCheckBox *graphicalCheckBox_;
    QCheckBox *debugCheckBox_;
    QSpinBox *memorySpinBox_;
    QSpinBox *cpusSpinBox_;
    QLineEdit *cmdlineEdit_;

    QTextCharFormat defaultFormat_;
    QTextCharFormat currentFormat_;
};

#endif
