#ifndef WORKSPACEWIDGET_H
#define WORKSPACEWIDGET_H

#include <QWidget>
#include <QLabel>
#include <QPushButton>
#include <QLineEdit>
#include <QProgressBar>
#include <QGroupBox>

class GrpcClient;

class WorkspaceWidget : public QWidget
{
    Q_OBJECT

public:
    explicit WorkspaceWidget(GrpcClient *grpcClient, QWidget *parent = nullptr);
    ~WorkspaceWidget();

public slots:
    void initWorkspace();
    void mountWorkspace();
    void unmountWorkspace();
    void refreshStatus();

private slots:
    void onInitProgress(int stage, int progress, const QString &message);
    void onStatusUpdated(bool mounted, const QString &mountPoint, qint64 sizeBytes, qint64 usedBytes, qint64 availBytes);

private:
    void setupUI();
    void connectSignals();
    QString formatBytes(qint64 bytes) const;

    GrpcClient *grpcClient_;

    QLabel *statusIconLabel_;
    QLabel *statusLabel_;
    QLabel *mountPointLabel_;
    QLabel *volumePathLabel_;
    QLabel *sizeLabel_;
    QLabel *usedLabel_;
    QLabel *availLabel_;

    QLineEdit *workspaceNameEdit_;
    QLineEdit *volumeSizeEdit_;

    QPushButton *initButton_;
    QPushButton *mountButton_;
    QPushButton *unmountButton_;
    QPushButton *refreshButton_;

    QProgressBar *progressBar_;
    QLabel *progressMessageLabel_;

    bool isMounted_;
    bool isInitializing_;
};

#endif
