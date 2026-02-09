#ifndef CODEEDITOR_H
#define CODEEDITOR_H

#include <QPlainTextEdit>
#include <QObject>

class LineNumberArea;
class SyntaxHighlighter;

class CodeEditor : public QPlainTextEdit
{
    Q_OBJECT

public:
    explicit CodeEditor(QWidget *parent = nullptr);
    ~CodeEditor() override;

    void lineNumberAreaPaintEvent(QPaintEvent *event);
    int lineNumberAreaWidth();

    void setFilePath(const QString &path);
    QString filePath() const { return filePath_; }
    
    bool isModified() const;
    void setModified(bool modified);
    
    bool saveFile();
    bool loadFile(const QString &path);

protected:
    void resizeEvent(QResizeEvent *event) override;
    void keyPressEvent(QKeyEvent *event) override;

private slots:
    void updateLineNumberAreaWidth(int newBlockCount);
    void highlightCurrentLine();
    void updateLineNumberArea(const QRect &rect, int dy);

private:
    void setupEditor();
    void detectLanguage();
    
    LineNumberArea *lineNumberArea_;
    SyntaxHighlighter *highlighter_;
    QString filePath_;
    bool modified_;
};

class LineNumberArea : public QWidget
{
public:
    explicit LineNumberArea(CodeEditor *editor) : QWidget(editor), codeEditor_(editor) {}

    QSize sizeHint() const override
    {
        return QSize(codeEditor_->lineNumberAreaWidth(), 0);
    }

protected:
    void paintEvent(QPaintEvent *event) override
    {
        codeEditor_->lineNumberAreaPaintEvent(event);
    }

private:
    CodeEditor *codeEditor_;
};

#endif // CODEEDITOR_H
