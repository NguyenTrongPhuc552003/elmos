#ifndef ARCHSELECTORWIDGET_H
#define ARCHSELECTORWIDGET_H

#include <QWidget>
#include <QComboBox>
#include <QLabel>
#include <QPushButton>
#include <QGroupBox>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QMap>
#include <QString>

class GrpcClient;

class ArchSelectorWidget : public QWidget
{
    Q_OBJECT

public:
    explicit ArchSelectorWidget(GrpcClient *client, QWidget *parent = nullptr);
    
    QString currentArchitecture() const;
    void setArchitecture(const QString &arch);
    
signals:
    void architectureChanged(const QString &arch);
    void toolchainStatusUpdated(const QString &arch, bool installed);

private slots:
    void onArchitectureSelected(int index);
    void refreshToolchainStatus();
    void installToolchainForCurrentArch();
    void showArchitectureDetails();

private:
    void setupUI();
    void populateArchitectures();
    void updateToolchainStatus(const QString &arch);
    void updateStatusDisplay();
    
    GrpcClient *grpcClient_;
    
    QComboBox *archComboBox_;
    QLabel *statusLabel_;
    QLabel *detailsLabel_;
    QPushButton *installButton_;
    QPushButton *refreshButton_;
    QPushButton *detailsButton_;
    
    QMap<QString, bool> toolchainStatus_;
    QMap<QString, QString> archDescriptions_;
    
    QString currentArch_;
};

#endif
