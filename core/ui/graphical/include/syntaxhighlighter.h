#ifndef SYNTAXHIGHLIGHTER_H
#define SYNTAXHIGHLIGHTER_H

#include <QSyntaxHighlighter>
#include <QTextCharFormat>
#include <QRegularExpression>

class SyntaxHighlighter : public QSyntaxHighlighter
{
    Q_OBJECT

public:
    enum Language {
        None,
        C,
        CPlusPlus,
        Rust,
        Makefile,
        Shell
    };

    explicit SyntaxHighlighter(QTextDocument *parent = nullptr);
    
    void setLanguage(Language lang);
    Language language() const { return language_; }

protected:
    void highlightBlock(const QString &text) override;

private:
    void setupCLanguage();
    void setupCPlusPlusLanguage();
    void setupRustLanguage();
    void setupMakefileLanguage();
    void setupShellLanguage();
    
    struct HighlightingRule
    {
        QRegularExpression pattern;
        QTextCharFormat format;
    };
    
    Language language_;
    QVector<HighlightingRule> highlightingRules_;
    
    QTextCharFormat keywordFormat_;
    QTextCharFormat typeFormat_;
    QTextCharFormat stringFormat_;
    QTextCharFormat numberFormat_;
    QTextCharFormat commentFormat_;
    QTextCharFormat preprocessorFormat_;
    QTextCharFormat functionFormat_;
    
    QRegularExpression commentStartExpression_;
    QRegularExpression commentEndExpression_;
};

#endif // SYNTAXHIGHLIGHTER_H
