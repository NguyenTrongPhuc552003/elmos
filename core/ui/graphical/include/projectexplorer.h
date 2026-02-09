#ifndef PROJECTEXPLORER_H
#define PROJECTEXPLORER_H

#include <QWidget>
#include <QTreeView>
#include <QFileSystemModel>
#include <QMenu>
#include <QAction>

class ProjectExplorer : public QWidget
{
    Q_OBJECT

public:
    explicit ProjectExplorer(QWidget *parent = nullptr);
    
    void setRootPath(const QString &path);
    QString rootPath() const;
    
    QString selectedFilePath() const;

signals:
    void fileDoubleClicked(const QString &filePath);
    void fileSelected(const QString &filePath);

private slots:
    void onDoubleClicked(const QModelIndex &index);
    void onCustomContextMenuRequested(const QPoint &pos);
    void newFile();
    void newFolder();
    void renameItem();
    void deleteItem();
    void refreshTree();

private:
    void setupUI();
    void createActions();
    void createContextMenu();
    
    QTreeView *treeView_;
    QFileSystemModel *fileSystemModel_;
    QMenu *contextMenu_;
    
    QAction *newFileAction_;
    QAction *newFolderAction_;
    QAction *renameAction_;
    QAction *deleteAction_;
    QAction *refreshAction_;
    
    QString rootPath_;
};

#endif // PROJECTEXPLORER_H
