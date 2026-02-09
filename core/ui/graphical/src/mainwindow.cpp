#include "mainwindow.h"
#include "codeeditor.h"
#include "projectwizard.h"
#include "projectexplorer.h"
#include "grpcclient.h"
#include "kernelbuildwidget.h"
#include "qemuconsolewidget.h"
#include "settingsdialog.h"
#include "workspacewidget.h"
#include "toolchainwidget.h"
#include "archselectorwidget.h"
#include <QAction>
#include <QFileDialog>
#include <QFileInfo>
#include <QMessageBox>
#include <QSettings>
#include <QCloseEvent>
#include <QLabel>
#include <QVBoxLayout>
#include <QHBoxLayout>
#include <QSplitter>
#include <QPushButton>
#include <QDialog>

MainWindow::MainWindow(QWidget *parent)
    : QMainWindow(parent)
    , editorTabs_(nullptr)
    , projectExplorerDock_(nullptr)
    , buildOutputDock_(nullptr)
    , qemuConsoleDock_(nullptr)
    , kernelBuildDock_(nullptr)
    , projectExplorer_(nullptr)
    , buildOutput_(nullptr)
    , qemuConsole_(nullptr)
    , kernelBuildWidget_(nullptr)
    , workspaceManager_(nullptr)
    , archSelector_(nullptr)
    , grpcClient_(nullptr)
    , fileMenu_(nullptr)
    , editMenu_(nullptr)
    , buildMenu_(nullptr)
    , debugMenu_(nullptr)
    , toolsMenu_(nullptr)
    , helpMenu_(nullptr)
    , fileToolBar_(nullptr)
    , buildToolBar_(nullptr)
    , debugToolBar_(nullptr)
    , workspaceMounted_(false)
{
    setWindowTitle("ELMOS - Embedded Linux Development IDE");
    resize(1400, 900);
    
    createActions();
    createMenus();
    createToolBars();
    createCentralWidget();
    createDockWidgets();
    createStatusBar();
    setupConnections();
    loadSettings();
}

MainWindow::~MainWindow()
{
    saveSettings();
}

void MainWindow::createActions()
{
    newProjectAction_ = new QAction(QIcon(":/icons/new-project.png"), tr("&New Project..."), this);
    newProjectAction_->setShortcut(QKeySequence::New);
    newProjectAction_->setStatusTip(tr("Create a new kernel module or user application project"));
    connect(newProjectAction_, &QAction::triggered, this, &MainWindow::newProject);
    
    openProjectAction_ = new QAction(QIcon(":/icons/open-folder.png"), tr("&Open Project..."), this);
    openProjectAction_->setShortcut(QKeySequence(tr("Ctrl+Shift+O")));
    openProjectAction_->setStatusTip(tr("Open an existing project"));
    connect(openProjectAction_, &QAction::triggered, this, &MainWindow::openProject);
    
    openFileAction_ = new QAction(QIcon(":/icons/open-file.png"), tr("Open &File..."), this);
    openFileAction_->setShortcut(QKeySequence::Open);
    openFileAction_->setStatusTip(tr("Open a file in the editor"));
    connect(openFileAction_, &QAction::triggered, this, &MainWindow::openFile);
    
    saveAction_ = new QAction(QIcon(":/icons/save.png"), tr("&Save"), this);
    saveAction_->setShortcut(QKeySequence::Save);
    saveAction_->setStatusTip(tr("Save the current file"));
    connect(saveAction_, &QAction::triggered, this, &MainWindow::saveFile);
    
    saveAsAction_ = new QAction(tr("Save &As..."), this);
    saveAsAction_->setShortcut(QKeySequence::SaveAs);
    saveAsAction_->setStatusTip(tr("Save the current file with a new name"));
    connect(saveAsAction_, &QAction::triggered, this, &MainWindow::saveFileAs);
    
    closeAction_ = new QAction(tr("&Close File"), this);
    closeAction_->setShortcut(QKeySequence::Close);
    closeAction_->setStatusTip(tr("Close the current file"));
    connect(closeAction_, &QAction::triggered, this, &MainWindow::closeFile);
    
    exitAction_ = new QAction(tr("E&xit"), this);
    exitAction_->setShortcut(QKeySequence::Quit);
    exitAction_->setStatusTip(tr("Exit the application"));
    connect(exitAction_, &QAction::triggered, this, &QWidget::close);
    
    buildKernelAction_ = new QAction(QIcon(":/icons/build-kernel.png"), tr("Build &Kernel"), this);
    buildKernelAction_->setShortcut(QKeySequence(tr("Ctrl+Shift+K")));
    buildKernelAction_->setStatusTip(tr("Build Linux kernel for current architecture"));
    connect(buildKernelAction_, &QAction::triggered, this, &MainWindow::buildKernel);
    
    buildModuleAction_ = new QAction(QIcon(":/icons/build-module.png"), tr("Build &Module"), this);
    buildModuleAction_->setShortcut(QKeySequence(tr("Ctrl+Shift+M")));
    buildModuleAction_->setStatusTip(tr("Build kernel module"));
    connect(buildModuleAction_, &QAction::triggered, this, &MainWindow::buildModule);
    
    buildUserAppAction_ = new QAction(QIcon(":/icons/build-app.png"), tr("Build &App"), this);
    buildUserAppAction_->setShortcut(QKeySequence(tr("Ctrl+Shift+B")));
    buildUserAppAction_->setStatusTip(tr("Build user application"));
    connect(buildUserAppAction_, &QAction::triggered, this, &MainWindow::buildUserApp);
    
    cleanAction_ = new QAction(QIcon(":/icons/clean.png"), tr("&Clean Build"), this);
    cleanAction_->setStatusTip(tr("Clean all build artifacts"));
    connect(cleanAction_, &QAction::triggered, this, &MainWindow::cleanBuild);
    
    runAction_ = new QAction(QIcon(":/icons/run.png"), tr("&Run QEMU"), this);
    runAction_->setShortcut(QKeySequence(tr("F5")));
    runAction_->setStatusTip(tr("Run kernel in QEMU emulator"));
    connect(runAction_, &QAction::triggered, this, &MainWindow::runQEMU);
    
    debugAction_ = new QAction(QIcon(":/icons/debug.png"), tr("&Debug QEMU"), this);
    debugAction_->setShortcut(QKeySequence(tr("F9")));
    debugAction_->setStatusTip(tr("Run QEMU with GDB debugging enabled"));
    connect(debugAction_, &QAction::triggered, this, &MainWindow::debugQEMU);
    
    stopAction_ = new QAction(QIcon(":/icons/stop.png"), tr("&Stop"), this);
    stopAction_->setShortcut(QKeySequence(tr("Shift+F5")));
    stopAction_->setStatusTip(tr("Stop QEMU emulator"));
    stopAction_->setEnabled(false);
    connect(stopAction_, &QAction::triggered, this, &MainWindow::stopQEMU);
    
    settingsAction_ = new QAction(QIcon(":/icons/settings.png"), tr("&Settings..."), this);
    settingsAction_->setShortcut(QKeySequence::Preferences);
    settingsAction_->setStatusTip(tr("Configure IDE settings"));
    connect(settingsAction_, &QAction::triggered, this, &MainWindow::showSettings);
    
    aboutAction_ = new QAction(tr("&About ELMOS"), this);
    aboutAction_->setStatusTip(tr("Show information about ELMOS IDE"));
    connect(aboutAction_, &QAction::triggered, this, &MainWindow::showAbout);
    
    initWorkspaceAction_ = new QAction(tr("&Initialize Workspace..."), this);
    initWorkspaceAction_->setStatusTip(tr("Initialize a new ELMOS workspace"));
    connect(initWorkspaceAction_, &QAction::triggered, this, &MainWindow::initWorkspace);
    
    mountWorkspaceAction_ = new QAction(tr("&Mount Workspace"), this);
    mountWorkspaceAction_->setStatusTip(tr("Mount the ELMOS workspace volume"));
    connect(mountWorkspaceAction_, &QAction::triggered, this, &MainWindow::mountWorkspace);
    
    unmountWorkspaceAction_ = new QAction(tr("&Unmount Workspace"), this);
    unmountWorkspaceAction_->setStatusTip(tr("Unmount the ELMOS workspace volume"));
    connect(unmountWorkspaceAction_, &QAction::triggered, this, &MainWindow::unmountWorkspace);
    
    manageWorkspaceAction_ = new QAction(tr("&Manage Workspace..."), this);
    manageWorkspaceAction_->setStatusTip(tr("Open workspace management dialog"));
    connect(manageWorkspaceAction_, &QAction::triggered, this, &MainWindow::manageWorkspace);
    
    manageToolchainsAction_ = new QAction(tr("&Manage Toolchains..."), this);
    manageToolchainsAction_->setStatusTip(tr("Open toolchain management dialog"));
    connect(manageToolchainsAction_, &QAction::triggered, this, &MainWindow::manageToolchains);
}

void MainWindow::createMenus()
{
    fileMenu_ = menuBar()->addMenu(tr("&File"));
    fileMenu_->addAction(newProjectAction_);
    fileMenu_->addAction(openProjectAction_);
    fileMenu_->addAction(openFileAction_);
    fileMenu_->addSeparator();
    fileMenu_->addAction(saveAction_);
    fileMenu_->addAction(saveAsAction_);
    fileMenu_->addSeparator();
    fileMenu_->addAction(closeAction_);
    fileMenu_->addSeparator();
    fileMenu_->addAction(exitAction_);
    
    editMenu_ = menuBar()->addMenu(tr("&Edit"));
    
    buildMenu_ = menuBar()->addMenu(tr("&Build"));
    buildMenu_->addAction(buildKernelAction_);
    buildMenu_->addAction(buildModuleAction_);
    buildMenu_->addAction(buildUserAppAction_);
    buildMenu_->addSeparator();
    buildMenu_->addAction(cleanAction_);
    
    debugMenu_ = menuBar()->addMenu(tr("&Debug"));
    debugMenu_->addAction(runAction_);
    debugMenu_->addAction(debugAction_);
    debugMenu_->addAction(stopAction_);
    
    toolsMenu_ = menuBar()->addMenu(tr("&Tools"));
    QMenu *workspaceMenu = toolsMenu_->addMenu(tr("&Workspace"));
    QAction *initWorkspaceAction = workspaceMenu->addAction(tr("&Initialize Workspace..."));
    connect(initWorkspaceAction, &QAction::triggered, this, &MainWindow::initWorkspace);
    QAction *mountWorkspaceAction = workspaceMenu->addAction(tr("&Mount Workspace"));
    connect(mountWorkspaceAction, &QAction::triggered, this, &MainWindow::mountWorkspace);
    QAction *unmountWorkspaceAction = workspaceMenu->addAction(tr("&Unmount Workspace"));
    connect(unmountWorkspaceAction, &QAction::triggered, this, &MainWindow::unmountWorkspace);
    workspaceMenu->addSeparator();
    QAction *manageWorkspaceAction = workspaceMenu->addAction(tr("&Manage Workspace..."));
    connect(manageWorkspaceAction, &QAction::triggered, this, &MainWindow::manageWorkspace);
    
    QMenu *toolchainMenu = toolsMenu_->addMenu(tr("Tool&chain"));
    QAction *manageToolchainsAction = toolchainMenu->addAction(tr("&Manage Toolchains..."));
    connect(manageToolchainsAction, &QAction::triggered, this, &MainWindow::manageToolchains);
    
    toolsMenu_->addSeparator();
    toolsMenu_->addAction(settingsAction_);
    
    helpMenu_ = menuBar()->addMenu(tr("&Help"));
    helpMenu_->addAction(aboutAction_);
}

void MainWindow::createToolBars()
{
    fileToolBar_ = addToolBar(tr("File"));
    fileToolBar_->setObjectName("FileToolBar");
    fileToolBar_->addAction(newProjectAction_);
    fileToolBar_->addAction(openProjectAction_);
    fileToolBar_->addAction(saveAction_);
    
    buildToolBar_ = addToolBar(tr("Build"));
    buildToolBar_->setObjectName("BuildToolBar");
    buildToolBar_->addAction(buildKernelAction_);
    buildToolBar_->addAction(buildModuleAction_);
    buildToolBar_->addAction(buildUserAppAction_);
    buildToolBar_->addAction(cleanAction_);
    
    debugToolBar_ = addToolBar(tr("Debug"));
    debugToolBar_->setObjectName("DebugToolBar");
    debugToolBar_->addAction(runAction_);
    debugToolBar_->addAction(debugAction_);
    debugToolBar_->addAction(stopAction_);
}

void MainWindow::createCentralWidget()
{
    editorTabs_ = new QTabWidget(this);
    editorTabs_->setTabsClosable(true);
    editorTabs_->setMovable(true);
    editorTabs_->setDocumentMode(true);
    
    connect(editorTabs_, &QTabWidget::tabCloseRequested, this, [this](int index) {
        QWidget *widget = editorTabs_->widget(index);
        editorTabs_->removeTab(index);
        widget->deleteLater();
    });
    
    setCentralWidget(editorTabs_);
}

void MainWindow::createDockWidgets()
{
    projectExplorerDock_ = new QDockWidget(tr("Project Explorer"), this);
    projectExplorerDock_->setObjectName("ProjectExplorerDock");
    projectExplorerDock_->setAllowedAreas(Qt::LeftDockWidgetArea | Qt::RightDockWidgetArea);
    
    projectExplorer_ = new ProjectExplorer(this);
    connect(projectExplorer_, &ProjectExplorer::fileDoubleClicked, 
            this, &MainWindow::openFileInEditor);
    projectExplorerDock_->setWidget(projectExplorer_);
    
    addDockWidget(Qt::LeftDockWidgetArea, projectExplorerDock_);
    
    buildOutputDock_ = new QDockWidget(tr("Build Output"), this);
    buildOutputDock_->setObjectName("BuildOutputDock");
    buildOutputDock_->setAllowedAreas(Qt::BottomDockWidgetArea | Qt::TopDockWidgetArea);
    
    QTextEdit *buildPlaceholder = new QTextEdit(this);
    buildPlaceholder->setReadOnly(true);
    buildPlaceholder->setPlaceholderText("Build output will appear here...");
    buildOutputDock_->setWidget(buildPlaceholder);
    
    addDockWidget(Qt::BottomDockWidgetArea, buildOutputDock_);
    
    qemuConsoleDock_ = new QDockWidget(tr("QEMU Console"), this);
    qemuConsoleDock_->setObjectName("QEMUConsoleDock");
    qemuConsoleDock_->setAllowedAreas(Qt::BottomDockWidgetArea | Qt::TopDockWidgetArea);
    
    qemuConsole_ = new QEMUConsoleWidget(grpcClient_, this);
    qemuConsoleDock_->setWidget(qemuConsole_);
    
    addDockWidget(Qt::BottomDockWidgetArea, qemuConsoleDock_);
    
    kernelBuildDock_ = new QDockWidget(tr("Kernel Build Status"), this);
    kernelBuildDock_->setObjectName("KernelBuildDock");
    kernelBuildDock_->setAllowedAreas(Qt::RightDockWidgetArea | Qt::LeftDockWidgetArea);
    
    grpcClient_ = new GrpcClient("unix:///tmp/elmos.sock", this);
    kernelBuildWidget_ = new KernelBuildWidget(grpcClient_, this);
    kernelBuildDock_->setWidget(kernelBuildWidget_);
    
    addDockWidget(Qt::RightDockWidgetArea, kernelBuildDock_);
    
    archSelector_ = new ArchSelectorWidget(grpcClient_, this);
    connect(archSelector_, &ArchSelectorWidget::architectureChanged,
            this, [this](const QString &arch) {
        statusBar()->showMessage(tr("Architecture changed to: %1").arg(arch), 3000);
    });
    
    tabifyDockWidget(buildOutputDock_, qemuConsoleDock_);
    buildOutputDock_->raise();
}

void MainWindow::createStatusBar()
{
    QLabel *statusLabel = new QLabel("Ready", this);
    statusBar()->addWidget(statusLabel);
    
    QPushButton *archButton = new QPushButton(this);
    archButton->setFlat(true);
    archButton->setText(tr("Architecture: arm"));
    archButton->setToolTip(tr("Click to change target architecture"));
    archButton->setCursor(Qt::PointingHandCursor);
    connect(archButton, &QPushButton::clicked, this, &MainWindow::showArchSelector);
    connect(archSelector_, &ArchSelectorWidget::architectureChanged,
            archButton, [archButton](const QString &arch) {
        archButton->setText(QString("Architecture: %1").arg(arch));
    });
    statusBar()->addPermanentWidget(archButton);
    
    QLabel *workspaceLabel = new QLabel("Workspace: Not mounted", this);
    statusBar()->addPermanentWidget(workspaceLabel);
}

void MainWindow::setupConnections()
{
}

void MainWindow::loadSettings()
{
    QSettings settings;
    
    restoreGeometry(settings.value("mainwindow/geometry").toByteArray());
    restoreState(settings.value("mainwindow/state").toByteArray());
    
    currentProjectPath_ = settings.value("project/lastPath").toString();
}

void MainWindow::saveSettings()
{
    QSettings settings;
    
    settings.setValue("mainwindow/geometry", saveGeometry());
    settings.setValue("mainwindow/state", saveState());
    settings.setValue("project/lastPath", currentProjectPath_);
}

void MainWindow::newProject()
{
    ProjectWizard wizard(currentProjectPath_.isEmpty() ? QDir::homePath() : currentProjectPath_, this);
    
    wizard.setField("workspacePath", currentProjectPath_.isEmpty() ? QDir::homePath() : currentProjectPath_);
    
    if (wizard.exec() == QDialog::Accepted) {
        QString projectName = wizard.projectName();
        QString projectPath = wizard.projectPath();
        bool isKernelModule = (wizard.projectType() == ProjectWizard::KernelModule);
        QString author = wizard.field("author").toString();
        QString license = wizard.field("license").toString();
        QString description = wizard.field("description").toString();
        
        QDir dir;
        if (!dir.mkpath(projectPath)) {
            QMessageBox::critical(this, tr("Error"),
                                tr("Could not create project directory: %1").arg(projectPath));
            return;
        }
        
        bool success = false;
        if (isKernelModule) {
            success = createKernelModuleFiles(projectPath, projectName, author, license, description);
        } else {
            success = createUserAppFiles(projectPath, projectName, author, license, description);
        }
        
        if (success) {
            currentProjectPath_ = QFileInfo(projectPath).dir().path();
            if (projectExplorer_) {
                projectExplorer_->setRootPath(projectPath);
            }
            statusBar()->showMessage(tr("Project created: %1").arg(projectName), 5000);
        } else {
            QMessageBox::critical(this, tr("Error"),
                                tr("Failed to create project files."));
        }
    }
}

bool MainWindow::createKernelModuleFiles(const QString &projectPath, const QString &projectName,
                                        const QString &author, const QString &license, const QString &description)
{
    QString sourceFile = projectPath + "/" + projectName + ".c";
    QFile file(sourceFile);
    if (!file.open(QIODevice::WriteOnly | QIODevice::Text))
        return false;
    
    QTextStream out(&file);
    out << "/*\n";
    out << " * " << projectName << ".c - " << description << "\n";
    out << " * Author: " << author << "\n";
    out << " * License: " << license << "\n";
    out << " */\n\n";
    out << "#include <linux/module.h>\n";
    out << "#include <linux/kernel.h>\n";
    out << "#include <linux/init.h>\n\n";
    out << "MODULE_LICENSE(\"" << license << "\");\n";
    out << "MODULE_AUTHOR(\"" << author << "\");\n";
    out << "MODULE_DESCRIPTION(\"" << description << "\");\n\n";
    out << "static int __init " << projectName << "_init(void)\n";
    out << "{\n";
    out << "    pr_info(\"" << projectName << ": module loaded\\n\");\n";
    out << "    return 0;\n";
    out << "}\n\n";
    out << "static void __exit " << projectName << "_exit(void)\n";
    out << "{\n";
    out << "    pr_info(\"" << projectName << ": module unloaded\\n\");\n";
    out << "}\n\n";
    out << "module_init(" << projectName << "_init);\n";
    out << "module_exit(" << projectName << "_exit);\n";
    file.close();
    
    QString makefile = projectPath + "/Makefile";
    QFile makeFile(makefile);
    if (!makeFile.open(QIODevice::WriteOnly | QIODevice::Text))
        return false;
    
    QTextStream makeOut(&makeFile);
    makeOut << "obj-m += " << projectName << ".o\n\n";
    makeOut << "KDIR := /path/to/kernel/source\n\n";
    makeOut << "all:\n";
    makeOut << "\tmake -C $(KDIR) M=$(PWD) modules\n\n";
    makeOut << "clean:\n";
    makeOut << "\tmake -C $(KDIR) M=$(PWD) clean\n";
    makeFile.close();
    
    createReadme(projectPath, projectName, "Kernel Module", description, author, license);
    
    return true;
}

bool MainWindow::createUserAppFiles(const QString &projectPath, const QString &projectName,
                                   const QString &author, const QString &license, const QString &description)
{
    QString sourceFile = projectPath + "/main.c";
    QFile file(sourceFile);
    if (!file.open(QIODevice::WriteOnly | QIODevice::Text))
        return false;
    
    QTextStream out(&file);
    out << "/*\n";
    out << " * " << projectName << " - " << description << "\n";
    out << " * Author: " << author << "\n";
    out << " * License: " << license << "\n";
    out << " */\n\n";
    out << "#include <stdio.h>\n";
    out << "#include <stdlib.h>\n\n";
    out << "int main(int argc, char *argv[])\n";
    out << "{\n";
    out << "    printf(\"Hello from " << projectName << "!\\n\");\n";
    out << "    return 0;\n";
    out << "}\n";
    file.close();
    
    QString makefile = projectPath + "/Makefile";
    QFile makeFile(makefile);
    if (!makeFile.open(QIODevice::WriteOnly | QIODevice::Text))
        return false;
    
    QTextStream makeOut(&makeFile);
    makeOut << "CC := $(CROSS_COMPILE)gcc\n";
    makeOut << "CFLAGS := -Wall -O2\n";
    makeOut << "TARGET := " << projectName << "\n\n";
    makeOut << "all: $(TARGET)\n\n";
    makeOut << "$(TARGET): main.c\n";
    makeOut << "\t$(CC) $(CFLAGS) -o $@ $<\n\n";
    makeOut << "clean:\n";
    makeOut << "\trm -f $(TARGET)\n";
    makeFile.close();
    
    createReadme(projectPath, projectName, "User Application", description, author, license);
    
    return true;
}

void MainWindow::createReadme(const QString &projectPath, const QString &projectName,
                             const QString &projectType, const QString &description,
                             const QString &author, const QString &license)
{
    QString readme = projectPath + "/README.md";
    QFile readmeFile(readme);
    if (readmeFile.open(QIODevice::WriteOnly | QIODevice::Text)) {
        QTextStream out(&readmeFile);
        out << "# " << projectName << "\n\n";
        out << description << "\n\n";
        out << "## Project Information\n\n";
        out << "- **Type:** " << projectType << "\n";
        out << "- **Author:** " << author << "\n";
        out << "- **License:** " << license << "\n\n";
        out << "## Building\n\n";
        out << "```bash\n";
        out << "make\n";
        out << "```\n\n";
        out << "## Usage\n\n";
        out << "TODO: Add usage instructions\n";
        readmeFile.close();
    }
}

void MainWindow::openProject()
{
    QString dir = QFileDialog::getExistingDirectory(
        this,
        tr("Open Project Directory"),
        currentProjectPath_.isEmpty() ? QDir::homePath() : currentProjectPath_,
        QFileDialog::ShowDirsOnly | QFileDialog::DontResolveSymlinks
    );
    
    if (!dir.isEmpty()) {
        currentProjectPath_ = dir;
        if (projectExplorer_) {
            projectExplorer_->setRootPath(dir);
        }
        statusBar()->showMessage(tr("Opened project: %1").arg(dir), 3000);
    }
}

void MainWindow::openFile()
{
    QString fileName = QFileDialog::getOpenFileName(
        this,
        tr("Open File"),
        currentProjectPath_.isEmpty() ? QDir::homePath() : currentProjectPath_,
        tr("All Files (*);;C Files (*.c *.h);;C++ Files (*.cpp *.hpp *.cc *.cxx);;Rust Files (*.rs);;Makefiles (Makefile* *.mk)")
    );
    
    if (!fileName.isEmpty()) {
        openFileInEditor(fileName);
    }
}

void MainWindow::openFileInEditor(const QString &filePath)
{
    for (int i = 0; i < editorTabs_->count(); ++i) {
        CodeEditor *editor = qobject_cast<CodeEditor*>(editorTabs_->widget(i));
        if (editor && editor->filePath() == filePath) {
            editorTabs_->setCurrentIndex(i);
            return;
        }
    }
    
    CodeEditor *editor = new CodeEditor(this);
    if (editor->loadFile(filePath)) {
        QFileInfo fileInfo(filePath);
        int index = editorTabs_->addTab(editor, fileInfo.fileName());
        editorTabs_->setCurrentIndex(index);
        statusBar()->showMessage(tr("Opened: %1").arg(filePath), 3000);
    } else {
        delete editor;
        QMessageBox::warning(this, tr("Open File"), 
                           tr("Could not open file: %1").arg(filePath));
    }
}

CodeEditor* MainWindow::currentEditor()
{
    return qobject_cast<CodeEditor*>(editorTabs_->currentWidget());
}

void MainWindow::saveFile()
{
    CodeEditor *editor = currentEditor();
    if (editor) {
        if (editor->saveFile()) {
            statusBar()->showMessage(tr("File saved"), 2000);
        } else {
            QMessageBox::warning(this, tr("Save File"),
                               tr("Could not save file: %1").arg(editor->filePath()));
        }
    }
}

void MainWindow::saveFileAs()
{
    QString fileName = QFileDialog::getSaveFileName(
        this,
        tr("Save File As"),
        currentProjectPath_,
        tr("All Files (*)")
    );
    
    if (!fileName.isEmpty()) {
        statusBar()->showMessage(tr("File saved as: %1").arg(fileName), 3000);
    }
}

void MainWindow::closeFile()
{
    int currentIndex = editorTabs_->currentIndex();
    if (currentIndex >= 0) {
        editorTabs_->removeTab(currentIndex);
    }
}

void MainWindow::buildKernel()
{
    statusBar()->showMessage(tr("Building kernel..."), 0);
    QMessageBox::information(this, tr("Build Kernel"),
                            tr("Kernel build will be implemented with gRPC streaming.\n\n"
                               "Will show:\n"
                               "- Real-time compilation progress\n"
                               "- Build logs with color coding\n"
                               "- Error navigation"));
}

void MainWindow::buildModule()
{
    statusBar()->showMessage(tr("Building module..."), 0);
    QMessageBox::information(this, tr("Build Module"),
                            tr("Module build will compile current kernel module project."));
}

void MainWindow::buildUserApp()
{
    statusBar()->showMessage(tr("Building user application..."), 0);
    QMessageBox::information(this, tr("Build User App"),
                            tr("User application build with cross-compiler."));
}

void MainWindow::cleanBuild()
{
    statusBar()->showMessage(tr("Cleaning build artifacts..."), 2000);
}

void MainWindow::runQEMU()
{
    statusBar()->showMessage(tr("Starting QEMU..."), 0);
    stopAction_->setEnabled(true);
    runAction_->setEnabled(false);
    
    QMessageBox::information(this, tr("Run QEMU"),
                            tr("QEMU will start with current kernel.\n\n"
                               "Console output will appear in QEMU Console tab."));
}

void MainWindow::debugQEMU()
{
    statusBar()->showMessage(tr("Starting QEMU with GDB..."), 0);
    stopAction_->setEnabled(true);
    debugAction_->setEnabled(false);
    
    QMessageBox::information(this, tr("Debug QEMU"),
                            tr("QEMU will start with GDB server on port 1234.\n\n"
                               "Connect with: gdb -ex 'target remote :1234'"));
}

void MainWindow::stopQEMU()
{
    statusBar()->showMessage(tr("Stopping QEMU..."), 2000);
    stopAction_->setEnabled(false);
    runAction_->setEnabled(true);
    debugAction_->setEnabled(true);
}

void MainWindow::initWorkspace()
{
    QMessageBox::information(this, tr("Initialize Workspace"),
                            tr("Workspace initialization will:\n\n"
                               "1. Create DMG volume (macOS) or directory (Linux)\n"
                               "2. Set up kernel sources\n"
                               "3. Configure toolchain paths\n"
                               "4. Prepare build environment"));
}

void MainWindow::mountWorkspace()
{
    statusBar()->showMessage(tr("Mounting workspace..."), 2000);
    workspaceMounted_ = true;
}

void MainWindow::unmountWorkspace()
{
    statusBar()->showMessage(tr("Unmounting workspace..."), 2000);
    workspaceMounted_ = false;
}

void MainWindow::manageWorkspace()
{
    QDialog *dialog = new QDialog(this);
    dialog->setWindowTitle(tr("Workspace Manager"));
    dialog->setMinimumSize(600, 400);
    
    QVBoxLayout *layout = new QVBoxLayout(dialog);
    WorkspaceWidget *workspaceWidget = new WorkspaceWidget(grpcClient_, dialog);
    layout->addWidget(workspaceWidget);
    
    QPushButton *closeButton = new QPushButton(tr("Close"), dialog);
    connect(closeButton, &QPushButton::clicked, dialog, &QDialog::accept);
    layout->addWidget(closeButton);
    
    dialog->exec();
    delete dialog;
}

void MainWindow::manageToolchains()
{
    QDialog *dialog = new QDialog(this);
    dialog->setWindowTitle(tr("Toolchain Manager"));
    dialog->setMinimumSize(700, 500);
    
    QVBoxLayout *layout = new QVBoxLayout(dialog);
    ToolchainWidget *toolchainWidget = new ToolchainWidget(grpcClient_, dialog);
    layout->addWidget(toolchainWidget);
    
    QPushButton *closeButton = new QPushButton(tr("Close"), dialog);
    connect(closeButton, &QPushButton::clicked, dialog, &QDialog::accept);
    layout->addWidget(closeButton);
    
    dialog->exec();
    delete dialog;
}

void MainWindow::showArchSelector()
{
    QDialog *dialog = new QDialog(this);
    dialog->setWindowTitle(tr("Architecture Selector"));
    dialog->setMinimumSize(500, 350);
    
    QVBoxLayout *layout = new QVBoxLayout(dialog);
    
    ArchSelectorWidget *archWidget = new ArchSelectorWidget(grpcClient_, dialog);
    archWidget->setArchitecture(archSelector_->currentArchitecture());
    
    connect(archWidget, &ArchSelectorWidget::architectureChanged,
            this, [this, archWidget](const QString &arch) {
        archSelector_->setArchitecture(arch);
    });
    
    layout->addWidget(archWidget);
    
    QPushButton *closeButton = new QPushButton(tr("Close"), dialog);
    connect(closeButton, &QPushButton::clicked, dialog, &QDialog::accept);
    layout->addWidget(closeButton);
    
    dialog->exec();
    delete dialog;
}

void MainWindow::showSettings()
{
    SettingsDialog dialog(this);
    if (dialog.exec() == QDialog::Accepted) {
        loadSettings();
        statusBar()->showMessage(tr("Settings updated"), 2000);
    }
}

void MainWindow::showAbout()
{
    QMessageBox::about(this, tr("About ELMOS IDE"),
                      tr("<h2>ELMOS IDE</h2>"
                         "<p><b>Embedded Linux on MacOS - Development Environment</b></p>"
                         "<p>Version: 1.0.0</p>"
                         "<p>A modern Qt 6 IDE for embedded Linux kernel and application development.</p>"
                         "<p><b>Features:</b></p>"
                         "<ul>"
                         "<li>Cross-platform kernel building</li>"
                         "<li>Integrated QEMU emulation</li>"
                         "<li>Real-time build monitoring via gRPC</li>"
                         "<li>Syntax-highlighted code editor</li>"
                         "<li>Project templates for modules and apps</li>"
                         "</ul>"
                         "<p>Built with Qt 6 and gRPC streaming</p>"));
}
