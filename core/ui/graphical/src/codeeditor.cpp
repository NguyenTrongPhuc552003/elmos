#include "codeeditor.h"
#include "syntaxhighlighter.h"
#include <QPainter>
#include <QTextBlock>
#include <QFileInfo>
#include <QFile>
#include <QTextStream>
#include <QKeyEvent>

CodeEditor::CodeEditor(QWidget *parent)
    : QPlainTextEdit(parent)
    , lineNumberArea_(new LineNumberArea(this))
    , highlighter_(new SyntaxHighlighter(document()))
    , modified_(false)
{
    setupEditor();
    
    connect(this, &CodeEditor::blockCountChanged, this, &CodeEditor::updateLineNumberAreaWidth);
    connect(this, &CodeEditor::updateRequest, this, &CodeEditor::updateLineNumberArea);
    connect(this, &CodeEditor::cursorPositionChanged, this, &CodeEditor::highlightCurrentLine);
    connect(this, &CodeEditor::textChanged, this, [this]() { setModified(true); });
    
    updateLineNumberAreaWidth(0);
    highlightCurrentLine();
}

CodeEditor::~CodeEditor()
{
}

void CodeEditor::setupEditor()
{
    QFont font;
    font.setFamily("Monaco");
    font.setPointSize(11);
    font.setFixedPitch(true);
    setFont(font);
    
    setTabStopDistance(fontMetrics().horizontalAdvance(' ') * 4);
    
    setLineWrapMode(QPlainTextEdit::NoWrap);
    
    QPalette p = palette();
    p.setColor(QPalette::Base, QColor(30, 30, 30));
    p.setColor(QPalette::Text, QColor(220, 220, 220));
    setPalette(p);
}

int CodeEditor::lineNumberAreaWidth()
{
    int digits = 1;
    int max = qMax(1, blockCount());
    while (max >= 10) {
        max /= 10;
        ++digits;
    }
    
    int space = 10 + fontMetrics().horizontalAdvance(QLatin1Char('9')) * digits;
    return space;
}

void CodeEditor::updateLineNumberAreaWidth(int)
{
    setViewportMargins(lineNumberAreaWidth(), 0, 0, 0);
}

void CodeEditor::updateLineNumberArea(const QRect &rect, int dy)
{
    if (dy)
        lineNumberArea_->scroll(0, dy);
    else
        lineNumberArea_->update(0, rect.y(), lineNumberArea_->width(), rect.height());
    
    if (rect.contains(viewport()->rect()))
        updateLineNumberAreaWidth(0);
}

void CodeEditor::resizeEvent(QResizeEvent *event)
{
    QPlainTextEdit::resizeEvent(event);
    
    QRect cr = contentsRect();
    lineNumberArea_->setGeometry(QRect(cr.left(), cr.top(), lineNumberAreaWidth(), cr.height()));
}

void CodeEditor::highlightCurrentLine()
{
    QList<QTextEdit::ExtraSelection> extraSelections;
    
    if (!isReadOnly()) {
        QTextEdit::ExtraSelection selection;
        
        QColor lineColor = QColor(50, 50, 50);
        
        selection.format.setBackground(lineColor);
        selection.format.setProperty(QTextFormat::FullWidthSelection, true);
        selection.cursor = textCursor();
        selection.cursor.clearSelection();
        extraSelections.append(selection);
    }
    
    setExtraSelections(extraSelections);
}

void CodeEditor::lineNumberAreaPaintEvent(QPaintEvent *event)
{
    QPainter painter(lineNumberArea_);
    painter.fillRect(event->rect(), QColor(40, 40, 40));
    
    QTextBlock block = firstVisibleBlock();
    int blockNumber = block.blockNumber();
    int top = qRound(blockBoundingGeometry(block).translated(contentOffset()).top());
    int bottom = top + qRound(blockBoundingRect(block).height());
    
    while (block.isValid() && top <= event->rect().bottom()) {
        if (block.isVisible() && bottom >= event->rect().top()) {
            QString number = QString::number(blockNumber + 1);
            painter.setPen(QColor(150, 150, 150));
            painter.drawText(0, top, lineNumberArea_->width() - 5, fontMetrics().height(),
                           Qt::AlignRight, number);
        }
        
        block = block.next();
        top = bottom;
        bottom = top + qRound(blockBoundingRect(block).height());
        ++blockNumber;
    }
}

void CodeEditor::keyPressEvent(QKeyEvent *event)
{
    if (event->key() == Qt::Key_Tab) {
        insertPlainText("    ");
        return;
    }
    
    if (event->key() == Qt::Key_Return || event->key() == Qt::Key_Enter) {
        QTextCursor cursor = textCursor();
        QString currentLine = cursor.block().text();
        
        int indent = 0;
        for (QChar c : currentLine) {
            if (c == ' ')
                indent++;
            else if (c == '\t')
                indent += 4;
            else
                break;
        }
        
        QPlainTextEdit::keyPressEvent(event);
        
        cursor = textCursor();
        QString indentation = QString(indent, ' ');
        cursor.insertText(indentation);
        
        return;
    }
    
    QPlainTextEdit::keyPressEvent(event);
}

void CodeEditor::setFilePath(const QString &path)
{
    filePath_ = path;
    detectLanguage();
}

void CodeEditor::detectLanguage()
{
    if (filePath_.isEmpty())
        return;
    
    QFileInfo fileInfo(filePath_);
    QString suffix = fileInfo.suffix().toLower();
    
    if (suffix == "c" || suffix == "h") {
        highlighter_->setLanguage(SyntaxHighlighter::C);
    } else if (suffix == "cpp" || suffix == "cc" || suffix == "cxx" || 
               suffix == "hpp" || suffix == "hh" || suffix == "hxx") {
        highlighter_->setLanguage(SyntaxHighlighter::CPlusPlus);
    } else if (suffix == "rs") {
        highlighter_->setLanguage(SyntaxHighlighter::Rust);
    } else if (fileInfo.fileName().toLower().startsWith("makefile") || suffix == "mk") {
        highlighter_->setLanguage(SyntaxHighlighter::Makefile);
    } else if (suffix == "sh" || suffix == "bash" || suffix == "zsh") {
        highlighter_->setLanguage(SyntaxHighlighter::Shell);
    } else {
        highlighter_->setLanguage(SyntaxHighlighter::None);
    }
}

bool CodeEditor::isModified() const
{
    return modified_;
}

void CodeEditor::setModified(bool modified)
{
    if (modified_ != modified) {
        modified_ = modified;
        document()->setModified(modified);
    }
}

bool CodeEditor::saveFile()
{
    if (filePath_.isEmpty())
        return false;
    
    QFile file(filePath_);
    if (!file.open(QIODevice::WriteOnly | QIODevice::Text))
        return false;
    
    QTextStream out(&file);
    out << toPlainText();
    
    setModified(false);
    return true;
}

bool CodeEditor::loadFile(const QString &path)
{
    QFile file(path);
    if (!file.open(QIODevice::ReadOnly | QIODevice::Text))
        return false;
    
    QTextStream in(&file);
    setPlainText(in.readAll());
    
    setFilePath(path);
    setModified(false);
    
    return true;
}
