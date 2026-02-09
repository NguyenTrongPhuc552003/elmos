#include "projectexplorer.h"
#include <QVBoxLayout>
#include <QHeaderView>
#include <QInputDialog>
#include <QMessageBox>
#include <QFile>
#include <QDir>

ProjectExplorer::ProjectExplorer(QWidget *parent)
    : QWidget(parent)
    , treeView_(new QTreeView(this))
    , fileSystemModel_(new QFileSystemModel(this))
{
    setupUI();
    createActions();
    createContextMenu();
}

void ProjectExplorer::setupUI()
{
    fileSystemModel_->setReadOnly(false);
    fileSystemModel_->setNameFilters({"*"});
    fileSystemModel_->setNameFilterDisables(false);
    
    treeView_->setModel(fileSystemModel_);
    treeView_->setColumnWidth(0, 250);
    treeView_->setAlternatingRowColors(true);
    treeView_->setAnimated(true);
    treeView_->setIndentation(20);
    treeView_->setSortingEnabled(true);
    treeView_->setContextMenuPolicy(Qt::CustomContextMenu);
    
    treeView_->header()->setSectionResizeMode(0, QHeaderView::Stretch);
    treeView_->hideColumn(1);
    treeView_->hideColumn(2);
    treeView_->hideColumn(3);
    
    QPalette p = treeView_->palette();
    p.setColor(QPalette::Base, QColor(40, 40, 40));
    p.setColor(QPalette::Text, QColor(220, 220, 220));
    p.setColor(QPalette::AlternateBase, QColor(45, 45, 45));
    treeView_->setPalette(p);
    
    connect(treeView_, &QTreeView::doubleClicked, this, &ProjectExplorer::onDoubleClicked);
    connect(treeView_, &QTreeView::customContextMenuRequested, 
            this, &ProjectExplorer::onCustomContextMenuRequested);
    
    QVBoxLayout *layout = new QVBoxLayout(this);
    layout->setContentsMargins(0, 0, 0, 0);
    layout->addWidget(treeView_);
    
    setLayout(layout);
}

void ProjectExplorer::createActions()
{
    newFileAction_ = new QAction(QIcon(":/icons/new-file.png"), tr("New File..."), this);
    newFileAction_->setShortcut(QKeySequence(tr("Ctrl+N")));
    connect(newFileAction_, &QAction::triggered, this, &ProjectExplorer::newFile);
    
    newFolderAction_ = new QAction(QIcon(":/icons/new-folder.png"), tr("New Folder..."), this);
    connect(newFolderAction_, &QAction::triggered, this, &ProjectExplorer::newFolder);
    
    renameAction_ = new QAction(tr("Rename..."), this);
    renameAction_->setShortcut(QKeySequence(tr("F2")));
    connect(renameAction_, &QAction::triggered, this, &ProjectExplorer::renameItem);
    
    deleteAction_ = new QAction(tr("Delete"), this);
    deleteAction_->setShortcut(QKeySequence::Delete);
    connect(deleteAction_, &QAction::triggered, this, &ProjectExplorer::deleteItem);
    
    refreshAction_ = new QAction(tr("Refresh"), this);
    refreshAction_->setShortcut(QKeySequence::Refresh);
    connect(refreshAction_, &QAction::triggered, this, &ProjectExplorer::refreshTree);
}

void ProjectExplorer::createContextMenu()
{
    contextMenu_ = new QMenu(this);
    contextMenu_->addAction(newFileAction_);
    contextMenu_->addAction(newFolderAction_);
    contextMenu_->addSeparator();
    contextMenu_->addAction(renameAction_);
    contextMenu_->addAction(deleteAction_);
    contextMenu_->addSeparator();
    contextMenu_->addAction(refreshAction_);
}

void ProjectExplorer::setRootPath(const QString &path)
{
    rootPath_ = path;
    QModelIndex rootIndex = fileSystemModel_->setRootPath(path);
    treeView_->setRootIndex(rootIndex);
    treeView_->expand(rootIndex);
}

QString ProjectExplorer::rootPath() const
{
    return rootPath_;
}

QString ProjectExplorer::selectedFilePath() const
{
    QModelIndex index = treeView_->currentIndex();
    if (index.isValid()) {
        return fileSystemModel_->filePath(index);
    }
    return QString();
}

void ProjectExplorer::onDoubleClicked(const QModelIndex &index)
{
    if (index.isValid()) {
        QString filePath = fileSystemModel_->filePath(index);
        QFileInfo fileInfo(filePath);
        
        if (fileInfo.isFile()) {
            emit fileDoubleClicked(filePath);
        }
    }
}

void ProjectExplorer::onCustomContextMenuRequested(const QPoint &pos)
{
    QModelIndex index = treeView_->indexAt(pos);
    
    renameAction_->setEnabled(index.isValid());
    deleteAction_->setEnabled(index.isValid());
    
    contextMenu_->exec(treeView_->mapToGlobal(pos));
}

void ProjectExplorer::newFile()
{
    QModelIndex index = treeView_->currentIndex();
    QString dirPath;
    
    if (index.isValid()) {
        QFileInfo fileInfo(fileSystemModel_->filePath(index));
        dirPath = fileInfo.isDir() ? fileInfo.filePath() : fileInfo.dir().path();
    } else {
        dirPath = rootPath_;
    }
    
    bool ok;
    QString fileName = QInputDialog::getText(this, tr("New File"),
                                            tr("File name:"), QLineEdit::Normal,
                                            tr("newfile.c"), &ok);
    
    if (ok && !fileName.isEmpty()) {
        QString filePath = dirPath + "/" + fileName;
        QFile file(filePath);
        
        if (file.exists()) {
            QMessageBox::warning(this, tr("File Exists"),
                               tr("A file with this name already exists."));
            return;
        }
        
        if (file.open(QIODevice::WriteOnly)) {
            file.close();
            QMessageBox::information(this, tr("Success"),
                                   tr("File created: %1").arg(fileName));
        } else {
            QMessageBox::critical(this, tr("Error"),
                                tr("Could not create file: %1").arg(file.errorString()));
        }
    }
}

void ProjectExplorer::newFolder()
{
    QModelIndex index = treeView_->currentIndex();
    QString dirPath;
    
    if (index.isValid()) {
        QFileInfo fileInfo(fileSystemModel_->filePath(index));
        dirPath = fileInfo.isDir() ? fileInfo.filePath() : fileInfo.dir().path();
    } else {
        dirPath = rootPath_;
    }
    
    bool ok;
    QString folderName = QInputDialog::getText(this, tr("New Folder"),
                                              tr("Folder name:"), QLineEdit::Normal,
                                              tr("newfolder"), &ok);
    
    if (ok && !folderName.isEmpty()) {
        QDir dir(dirPath);
        if (dir.mkdir(folderName)) {
            QMessageBox::information(this, tr("Success"),
                                   tr("Folder created: %1").arg(folderName));
        } else {
            QMessageBox::critical(this, tr("Error"),
                                tr("Could not create folder."));
        }
    }
}

void ProjectExplorer::renameItem()
{
    QModelIndex index = treeView_->currentIndex();
    if (!index.isValid())
        return;
    
    QString oldPath = fileSystemModel_->filePath(index);
    QFileInfo fileInfo(oldPath);
    QString oldName = fileInfo.fileName();
    
    bool ok;
    QString newName = QInputDialog::getText(this, tr("Rename"),
                                           tr("New name:"), QLineEdit::Normal,
                                           oldName, &ok);
    
    if (ok && !newName.isEmpty() && newName != oldName) {
        QString newPath = fileInfo.dir().filePath(newName);
        
        if (QFile::exists(newPath)) {
            QMessageBox::warning(this, tr("Rename Failed"),
                               tr("A file or folder with this name already exists."));
            return;
        }
        
        if (QFile::rename(oldPath, newPath)) {
            QMessageBox::information(this, tr("Success"),
                                   tr("Renamed to: %1").arg(newName));
        } else {
            QMessageBox::critical(this, tr("Error"),
                                tr("Could not rename file or folder."));
        }
    }
}

void ProjectExplorer::deleteItem()
{
    QModelIndex index = treeView_->currentIndex();
    if (!index.isValid())
        return;
    
    QString path = fileSystemModel_->filePath(index);
    QFileInfo fileInfo(path);
    
    QMessageBox::StandardButton reply = QMessageBox::question(
        this, tr("Confirm Delete"),
        tr("Are you sure you want to delete:\n%1").arg(fileInfo.fileName()),
        QMessageBox::Yes | QMessageBox::No
    );
    
    if (reply == QMessageBox::Yes) {
        bool success;
        if (fileInfo.isDir()) {
            QDir dir(path);
            success = dir.removeRecursively();
        } else {
            success = QFile::remove(path);
        }
        
        if (success) {
            QMessageBox::information(this, tr("Success"), tr("Deleted successfully."));
        } else {
            QMessageBox::critical(this, tr("Error"), tr("Could not delete item."));
        }
    }
}

void ProjectExplorer::refreshTree()
{
    if (!rootPath_.isEmpty()) {
        setRootPath(rootPath_);
    }
}
