#ifndef MAINWINDOW_H
#define MAINWINDOW_H

#include <QMainWindow>
#include <QTabWidget>
#include <QDockWidget>
#include <QTextEdit>
#include <QToolBar>
#include <QMenuBar>
#include <QStatusBar>
#include <QTreeView>
#include <QSplitter>
#include <memory>

class CodeEditor;
class ProjectExplorer;
class ProjectWizard;
class WorkspaceManager;
class KernelBuildWidget;
class QEMUConsoleWidget;
class BuildOutputWidget;
class GrpcClient;
class ArchSelectorWidget;

class MainWindow : public QMainWindow
{
    Q_OBJECT

public:
    explicit MainWindow(QWidget *parent = nullptr);
    ~MainWindow() override;

private slots:
    void newProject();
    void openProject();
    void openFile();
    void saveFile();
    void saveFileAs();
    void closeFile();
    
    void buildKernel();
    void buildModule();
    void buildUserApp();
    void cleanBuild();
    
    void runQEMU();
    void debugQEMU();
    void stopQEMU();
    
    void initWorkspace();
    void mountWorkspace();
    void unmountWorkspace();
    void manageToolchains();
    void manageWorkspace();
    void showArchSelector();
    
    void showSettings();
    void showAbout();

private:
    void createActions();
    void createMenus();
    void createToolBars();
    void createDockWidgets();
    void createCentralWidget();
    void createStatusBar();
    void setupConnections();
    void loadSettings();
    void saveSettings();
    
    void openFileInEditor(const QString &filePath);
    CodeEditor* currentEditor();
    
    bool createKernelModuleFiles(const QString &projectPath, const QString &projectName,
                                const QString &author, const QString &license, const QString &description);
    bool createUserAppFiles(const QString &projectPath, const QString &projectName,
                           const QString &author, const QString &license, const QString &description);
    void createReadme(const QString &projectPath, const QString &projectName,
                     const QString &projectType, const QString &description,
                     const QString &author, const QString &license);

    // Central editor area
    QTabWidget *editorTabs_;
    
    // Dock widgets
    QDockWidget *projectExplorerDock_;
    QDockWidget *buildOutputDock_;
    QDockWidget *qemuConsoleDock_;
    QDockWidget *kernelBuildDock_;
    
    // Widgets
    ProjectExplorer *projectExplorer_;
    BuildOutputWidget *buildOutput_;
    QEMUConsoleWidget *qemuConsole_;
    KernelBuildWidget *kernelBuildWidget_;
    WorkspaceManager *workspaceManager_;
    ArchSelectorWidget *archSelector_;
    
    // gRPC client
    GrpcClient *grpcClient_;
    
    // Menus
    QMenu *fileMenu_;
    QMenu *editMenu_;
    QMenu *buildMenu_;
    QMenu *debugMenu_;
    QMenu *toolsMenu_;
    QMenu *helpMenu_;
    
    // Toolbars
    QToolBar *fileToolBar_;
    QToolBar *buildToolBar_;
    QToolBar *debugToolBar_;
    
    // Actions
    QAction *newProjectAction_;
    QAction *openProjectAction_;
    QAction *openFileAction_;
    QAction *saveAction_;
    QAction *saveAsAction_;
    QAction *closeAction_;
    QAction *exitAction_;
    
    QAction *buildKernelAction_;
    QAction *buildModuleAction_;
    QAction *buildUserAppAction_;
    QAction *cleanAction_;
    
    QAction *runAction_;
    QAction *debugAction_;
    QAction *stopAction_;
    
    QAction *initWorkspaceAction_;
    QAction *mountWorkspaceAction_;
    QAction *unmountWorkspaceAction_;
    QAction *manageWorkspaceAction_;
    QAction *manageToolchainsAction_;
    
    QAction *settingsAction_;
    QAction *aboutAction_;
    
    QString currentProjectPath_;
    bool workspaceMounted_;
};

#endif // MAINWINDOW_H
