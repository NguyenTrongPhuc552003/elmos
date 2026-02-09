#include "projectwizard.h"
#include <QVBoxLayout>
#include <QFormLayout>
#include <QMessageBox>
#include <QFileInfo>
#include <QDir>
#include <QFile>
#include <QTextStream>

ProjectWizard::ProjectWizard(const QString &workspacePath, QWidget *parent)
    : QWizard(parent)
    , workspacePath_(workspacePath)
{
    setWindowTitle(tr("New Project Wizard"));
    setWizardStyle(QWizard::ModernStyle);
    setOption(QWizard::HaveHelpButton, false);
    
    setPage(Page_Intro, new IntroPage);
    setPage(Page_Type, new ProjectTypePage);
    setPage(Page_Details, new ProjectDetailsPage);
    setPage(Page_Summary, new SummaryPage);
    
    setStartId(Page_Intro);
    
    resize(600, 500);
}

QString ProjectWizard::projectName() const
{
    return field("projectName").toString();
}

QString ProjectWizard::projectPath() const
{
    QString name = projectName();
    return workspacePath_ + "/" + name;
}

ProjectWizard::ProjectType ProjectWizard::projectType() const
{
    return field("isModule").toBool() ? KernelModule : UserApplication;
}

QString ProjectWizard::description() const
{
    return field("description").toString();
}

IntroPage::IntroPage(QWidget *parent)
    : QWizardPage(parent)
{
    setTitle(tr("Welcome to Project Creation"));
    setSubTitle(tr("This wizard will help you create a new kernel module or user application project."));
    
    label_ = new QLabel(
        tr("<p>ELMOS supports two types of projects:</p>"
           "<ul>"
           "<li><b>Kernel Module</b>: A loadable kernel module (LKM) that runs in kernel space. "
           "Useful for device drivers, file systems, and kernel extensions.</li>"
           "<li><b>User Application</b>: A userspace program compiled with the cross-compiler. "
           "Runs in your custom embedded Linux environment.</li>"
           "</ul>"
           "<p>Click <b>Next</b> to continue.</p>")
    );
    label_->setWordWrap(true);
    
    QVBoxLayout *layout = new QVBoxLayout;
    layout->addWidget(label_);
    setLayout(layout);
}

ProjectTypePage::ProjectTypePage(QWidget *parent)
    : QWizardPage(parent)
{
    setTitle(tr("Select Project Type"));
    setSubTitle(tr("Choose the type of project you want to create."));
    
    moduleRadio_ = new QRadioButton(tr("Kernel Module (LKM)"));
    appRadio_ = new QRadioButton(tr("User Application"));
    
    moduleDescLabel_ = new QLabel(
        tr("<small>Creates a kernel module with Makefile, source template, and Kbuild infrastructure. "
           "The module will be built against the current kernel source.</small>")
    );
    moduleDescLabel_->setWordWrap(true);
    moduleDescLabel_->setIndent(20);
    moduleDescLabel_->setStyleSheet("color: #666;");
    
    appDescLabel_ = new QLabel(
        tr("<small>Creates a userspace application with Makefile and cross-compilation setup. "
           "The binary will run in your embedded Linux rootfs.</small>")
    );
    appDescLabel_->setWordWrap(true);
    appDescLabel_->setIndent(20);
    appDescLabel_->setStyleSheet("color: #666;");
    
    moduleRadio_->setChecked(true);
    
    registerField("isModule", moduleRadio_);
    
    QVBoxLayout *layout = new QVBoxLayout;
    layout->addWidget(moduleRadio_);
    layout->addWidget(moduleDescLabel_);
    layout->addSpacing(20);
    layout->addWidget(appRadio_);
    layout->addWidget(appDescLabel_);
    layout->addStretch();
    
    setLayout(layout);
}

ProjectDetailsPage::ProjectDetailsPage(QWidget *parent)
    : QWizardPage(parent)
{
    setTitle(tr("Project Details"));
    setSubTitle(tr("Enter information about your project."));
    
    nameEdit_ = new QLineEdit;
    nameEdit_->setPlaceholderText(tr("e.g., my_driver or hello_world"));
    connect(nameEdit_, &QLineEdit::textChanged, this, &ProjectDetailsPage::updateCompleteState);
    
    authorEdit_ = new QLineEdit;
    authorEdit_->setPlaceholderText(tr("Your name"));
    
    descriptionEdit_ = new QTextEdit;
    descriptionEdit_->setPlaceholderText(tr("Brief description of your project..."));
    descriptionEdit_->setMaximumHeight(80);
    
    licenseCombo_ = new QComboBox;
    licenseCombo_->addItems({
        tr("GPL-2.0"),
        tr("GPL-2.0-or-later"),
        tr("MIT"),
        tr("BSD-3-Clause"),
        tr("Apache-2.0")
    });
    
    pathLabel_ = new QLabel;
    pathLabel_->setStyleSheet("color: #888; font-style: italic;");
    
    QFormLayout *layout = new QFormLayout;
    layout->addRow(tr("Project &Name:"), nameEdit_);
    layout->addRow(tr("&Author:"), authorEdit_);
    layout->addRow(tr("&License:"), licenseCombo_);
    layout->addRow(tr("&Description:"), descriptionEdit_);
    layout->addRow(tr("Project Path:"), pathLabel_);
    
    setLayout(layout);
    
    registerField("projectName*", nameEdit_);
    registerField("author", authorEdit_);
    registerField("license", licenseCombo_, "currentText");
    registerField("description", descriptionEdit_, "plainText");
}

bool ProjectDetailsPage::validatePage()
{
    QString name = nameEdit_->text().trimmed();
    
    if (name.isEmpty()) {
        QMessageBox::warning(this, tr("Validation Error"), 
                           tr("Project name cannot be empty."));
        return false;
    }
    
    if (!name.contains(QRegularExpression("^[a-zA-Z0-9_-]+$"))) {
        QMessageBox::warning(this, tr("Validation Error"),
                           tr("Project name can only contain letters, numbers, underscores, and hyphens."));
        return false;
    }
    
    QString projectPath = wizard()->field("workspacePath").toString() + "/" + name;
    QDir dir(projectPath);
    if (dir.exists()) {
        QMessageBox::warning(this, tr("Validation Error"),
                           tr("A project with this name already exists.\n\nPath: %1").arg(projectPath));
        return false;
    }
    
    return true;
}

void ProjectDetailsPage::initializePage()
{
    updateCompleteState();
}

void ProjectDetailsPage::updateCompleteState()
{
    QString name = nameEdit_->text().trimmed();
    QString basePath = wizard()->field("workspacePath").toString();
    
    if (!basePath.isEmpty() && !name.isEmpty()) {
        pathLabel_->setText(basePath + "/" + name);
    } else {
        pathLabel_->setText(tr("<not set>"));
    }
    
    emit completeChanged();
}

SummaryPage::SummaryPage(QWidget *parent)
    : QWizardPage(parent)
{
    setTitle(tr("Summary"));
    setSubTitle(tr("Review your project settings before creation."));
    
    summaryLabel_ = new QLabel;
    summaryLabel_->setWordWrap(true);
    summaryLabel_->setTextFormat(Qt::RichText);
    
    QVBoxLayout *layout = new QVBoxLayout;
    layout->addWidget(summaryLabel_);
    layout->addStretch();
    
    setLayout(layout);
}

void SummaryPage::initializePage()
{
    bool isModule = field("isModule").toBool();
    QString name = field("projectName").toString();
    QString author = field("author").toString();
    QString license = field("license").toString();
    QString description = field("description").toString();
    QString projectPath = wizard()->property("projectPath").toString();
    
    QString typeStr = isModule ? tr("<b>Kernel Module</b>") : tr("<b>User Application</b>");
    
    QString summary = QString(
        "<p><b>Ready to create your project!</b></p>"
        "<table cellspacing='8'>"
        "<tr><td align='right'><b>Type:</b></td><td>%1</td></tr>"
        "<tr><td align='right'><b>Name:</b></td><td>%2</td></tr>"
        "<tr><td align='right'><b>Author:</b></td><td>%3</td></tr>"
        "<tr><td align='right'><b>License:</b></td><td>%4</td></tr>"
        "<tr><td align='right'><b>Description:</b></td><td>%5</td></tr>"
        "<tr><td align='right'><b>Path:</b></td><td><tt>%6</tt></td></tr>"
        "</table>"
        "<p>The following files will be created:</p>"
        "<ul>"
    ).arg(typeStr, name, author.isEmpty() ? tr("<i>not set</i>") : author, 
          license, description.isEmpty() ? tr("<i>none</i>") : description,
          projectPath);
    
    if (isModule) {
        summary += QString(
            "<li><tt>%1.c</tt> - Module source code</li>"
            "<li><tt>Makefile</tt> - Kernel module build configuration</li>"
            "<li><tt>README.md</tt> - Project documentation</li>"
        ).arg(name);
    } else {
        summary += QString(
            "<li><tt>main.c</tt> - Application entry point</li>"
            "<li><tt>Makefile</tt> - Cross-compilation configuration</li>"
            "<li><tt>README.md</tt> - Project documentation</li>"
        ).arg(name);
    }
    
    summary += "</ul>";
    
    summaryLabel_->setText(summary);
}
