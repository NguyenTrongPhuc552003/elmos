#ifndef PROJECTWIZARD_H
#define PROJECTWIZARD_H

#include <QWizard>
#include <QWizardPage>
#include <QLineEdit>
#include <QTextEdit>
#include <QRadioButton>
#include <QComboBox>
#include <QLabel>

class ProjectWizard : public QWizard
{
    Q_OBJECT

public:
    enum ProjectType {
        KernelModule,
        UserApplication
    };

    enum { Page_Intro, Page_Type, Page_Details, Page_Summary };

    explicit ProjectWizard(const QString &workspacePath, QWidget *parent = nullptr);
    
    QString projectName() const;
    QString projectPath() const;
    ProjectType projectType() const;
    QString description() const;

private:
    QString workspacePath_;
};

class IntroPage : public QWizardPage
{
    Q_OBJECT

public:
    explicit IntroPage(QWidget *parent = nullptr);

private:
    QLabel *label_;
};

class ProjectTypePage : public QWizardPage
{
    Q_OBJECT

public:
    explicit ProjectTypePage(QWidget *parent = nullptr);

private:
    QRadioButton *moduleRadio_;
    QRadioButton *appRadio_;
    QLabel *moduleDescLabel_;
    QLabel *appDescLabel_;
};

class ProjectDetailsPage : public QWizardPage
{
    Q_OBJECT

public:
    explicit ProjectDetailsPage(QWidget *parent = nullptr);
    
    bool validatePage() override;
    void initializePage() override;

private slots:
    void updateCompleteState();

private:
    QLineEdit *nameEdit_;
    QLineEdit *authorEdit_;
    QTextEdit *descriptionEdit_;
    QComboBox *licenseCombo_;
    QLabel *pathLabel_;
};

class SummaryPage : public QWizardPage
{
    Q_OBJECT

public:
    explicit SummaryPage(QWidget *parent = nullptr);
    
    void initializePage() override;

private:
    QLabel *summaryLabel_;
};

#endif // PROJECTWIZARD_H
