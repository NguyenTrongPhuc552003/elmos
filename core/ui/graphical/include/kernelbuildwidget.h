#ifndef KERNELBUILDWIDGET_H
#define KERNELBUILDWIDGET_H

#include <QWidget>
#include <QTextEdit>
#include <QProgressBar>
#include <QLabel>
#include <QPushButton>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QComboBox>
#include <QSpinBox>
#include <QCheckBox>

class GrpcClient;

class KernelBuildWidget : public QWidget
{
    Q_OBJECT

public:
    explicit KernelBuildWidget(GrpcClient *grpcClient, QWidget *parent = nullptr);
    ~KernelBuildWidget();

public slots:
    void startBuild();
    void stopBuild();
    void clearOutput();

private slots:
    void onBuildStageChanged(const QString &stageName, int progress, int currentFile, int totalFiles);
    void onBuildLogReceived(int level, const QString &message, qint64 timestamp);
    void onBuildErrorReceived(const QString &message, const QString &file, int line);
    void onBuildCompleted(bool success, qint64 durationMs, const QString &imagePath);
    void onErrorOccurred(const QString &error);

private:
    void setupUI();
    void connectSignals();
    QString formatTimestamp(qint64 timestampMs) const;
    QString logLevelToString(int level) const;
    QString logLevelToColor(int level) const;

    GrpcClient *grpcClient_;
    bool isBuilding_;

    // UI components
    QTextEdit *outputView_;
    QProgressBar *progressBar_;
    QLabel *statusLabel_;
    QLabel *stageLabel_;
    QPushButton *buildButton_;
    QPushButton *stopButton_;
    QPushButton *clearButton_;
    
    // Build configuration
    QComboBox *archCombo_;
    QSpinBox *jobsSpinBox_;
    QCheckBox *verboseCheckBox_;
    QComboBox *targetCombo_;
};

#endif
