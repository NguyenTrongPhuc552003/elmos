#ifndef TOOLCHAINWIDGET_H
#define TOOLCHAINWIDGET_H

#include <QWidget>
#include <QListWidget>
#include <QPushButton>
#include <QLabel>
#include <QProgressBar>
#include <QTextEdit>
#include <QGroupBox>

class GrpcClient;

class ToolchainWidget : public QWidget
{
    Q_OBJECT

public:
    explicit ToolchainWidget(GrpcClient *grpcClient, QWidget *parent = nullptr);
    ~ToolchainWidget();

public slots:
    void refreshTargets();
    void refreshStatus();
    void installCrosstoolNG();
    void buildSelectedToolchain();
    void cleanToolchain();

private slots:
    void onTargetSelected(QListWidgetItem *item);
    void onInstallProgress(int stage, int progress, const QString &message);
    void onBuildProgress(int progress, const QString &message);

private:
    void setupUI();
    void connectSignals();

    GrpcClient *grpcClient_;

    QListWidget *targetsListWidget_;
    QTextEdit *detailsTextEdit_;
    QLabel *statusLabel_;
    QProgressBar *progressBar_;
    QPushButton *installButton_;
    QPushButton *buildButton_;
    QPushButton *cleanButton_;
    QPushButton *refreshButton_;

    QLabel *selectedTargetLabel_;
    QLabel *gccVersionLabel_;
    QLabel *pathLabel_;
    QLabel *statusIconLabel_;

    QString selectedTarget_;
    bool isBuilding_;
};

#endif
