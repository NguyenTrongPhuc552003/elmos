#include "syntaxhighlighter.h"

SyntaxHighlighter::SyntaxHighlighter(QTextDocument *parent)
    : QSyntaxHighlighter(parent)
    , language_(None)
{
    keywordFormat_.setForeground(QColor(86, 156, 214));
    keywordFormat_.setFontWeight(QFont::Bold);
    
    typeFormat_.setForeground(QColor(78, 201, 176));
    
    stringFormat_.setForeground(QColor(206, 145, 120));
    
    numberFormat_.setForeground(QColor(181, 206, 168));
    
    commentFormat_.setForeground(QColor(106, 153, 85));
    commentFormat_.setFontItalic(true);
    
    preprocessorFormat_.setForeground(QColor(189, 99, 197));
    
    functionFormat_.setForeground(QColor(220, 220, 170));
}

void SyntaxHighlighter::setLanguage(Language lang)
{
    if (language_ == lang)
        return;
    
    language_ = lang;
    highlightingRules_.clear();
    
    switch (lang) {
    case C:
        setupCLanguage();
        break;
    case CPlusPlus:
        setupCPlusPlusLanguage();
        break;
    case Rust:
        setupRustLanguage();
        break;
    case Makefile:
        setupMakefileLanguage();
        break;
    case Shell:
        setupShellLanguage();
        break;
    default:
        break;
    }
    
    rehighlight();
}

void SyntaxHighlighter::setupCLanguage()
{
    QStringList keywordPatterns = {
        "\\bauto\\b", "\\bbreak\\b", "\\bcase\\b", "\\bchar\\b", "\\bconst\\b",
        "\\bcontinue\\b", "\\bdefault\\b", "\\bdo\\b", "\\bdouble\\b", "\\belse\\b",
        "\\benum\\b", "\\bextern\\b", "\\bfloat\\b", "\\bfor\\b", "\\bgoto\\b",
        "\\bif\\b", "\\binline\\b", "\\bint\\b", "\\blong\\b", "\\bregister\\b",
        "\\brestrict\\b", "\\breturn\\b", "\\bshort\\b", "\\bsigned\\b", "\\bsizeof\\b",
        "\\bstatic\\b", "\\bstruct\\b", "\\bswitch\\b", "\\btypedef\\b", "\\bunion\\b",
        "\\bunsigned\\b", "\\bvoid\\b", "\\bvolatile\\b", "\\bwhile\\b"
    };
    
    for (const QString &pattern : keywordPatterns) {
        HighlightingRule rule;
        rule.pattern = QRegularExpression(pattern);
        rule.format = keywordFormat_;
        highlightingRules_.append(rule);
    }
    
    QStringList typePatterns = {
        "\\bint8_t\\b", "\\bint16_t\\b", "\\bint32_t\\b", "\\bint64_t\\b",
        "\\buint8_t\\b", "\\buint16_t\\b", "\\buint32_t\\b", "\\buint64_t\\b",
        "\\bsize_t\\b", "\\bssize_t\\b", "\\bbool\\b"
    };
    
    for (const QString &pattern : typePatterns) {
        HighlightingRule rule;
        rule.pattern = QRegularExpression(pattern);
        rule.format = typeFormat_;
        highlightingRules_.append(rule);
    }
    
    HighlightingRule rule;
    
    rule.pattern = QRegularExpression("\"([^\"\\\\]|\\\\.)*\"");
    rule.format = stringFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("'([^'\\\\]|\\\\.)*'");
    rule.format = stringFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("\\b[0-9]+\\b");
    rule.format = numberFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("//[^\n]*");
    rule.format = commentFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("#\\s*\\w+");
    rule.format = preprocessorFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("\\b[A-Za-z_][A-Za-z0-9_]*(?=\\s*\\()");
    rule.format = functionFormat_;
    highlightingRules_.append(rule);
    
    commentStartExpression_ = QRegularExpression("/\\*");
    commentEndExpression_ = QRegularExpression("\\*/");
}

void SyntaxHighlighter::setupCPlusPlusLanguage()
{
    setupCLanguage();
    
    QStringList cppKeywords = {
        "\\bclass\\b", "\\bnamespace\\b", "\\bpublic\\b", "\\bprivate\\b", "\\bprotected\\b",
        "\\bvirtual\\b", "\\boverride\\b", "\\bfinal\\b", "\\bexplicit\\b", "\\bconstexpr\\b",
        "\\bnoexcept\\b", "\\btemplate\\b", "\\btypename\\b", "\\boperator\\b", "\\bnew\\b",
        "\\bdelete\\b", "\\btry\\b", "\\bcatch\\b", "\\bthrow\\b", "\\busing\\b"
    };
    
    for (const QString &pattern : cppKeywords) {
        HighlightingRule rule;
        rule.pattern = QRegularExpression(pattern);
        rule.format = keywordFormat_;
        highlightingRules_.append(rule);
    }
}

void SyntaxHighlighter::setupRustLanguage()
{
    QStringList keywordPatterns = {
        "\\bas\\b", "\\bbreak\\b", "\\bconst\\b", "\\bcontinue\\b", "\\bcrate\\b",
        "\\belse\\b", "\\benum\\b", "\\bextern\\b", "\\bfn\\b", "\\bfor\\b",
        "\\bif\\b", "\\bimpl\\b", "\\bin\\b", "\\blet\\b", "\\bloop\\b",
        "\\bmatch\\b", "\\bmod\\b", "\\bmove\\b", "\\bmut\\b", "\\bpub\\b",
        "\\bref\\b", "\\breturn\\b", "\\bself\\b", "\\bSelf\\b", "\\bstatic\\b",
        "\\bstruct\\b", "\\bsuper\\b", "\\btrait\\b", "\\btype\\b", "\\bunsafe\\b",
        "\\buse\\b", "\\bwhere\\b", "\\bwhile\\b", "\\basync\\b", "\\bawait\\b"
    };
    
    for (const QString &pattern : keywordPatterns) {
        HighlightingRule rule;
        rule.pattern = QRegularExpression(pattern);
        rule.format = keywordFormat_;
        highlightingRules_.append(rule);
    }
    
    QStringList typePatterns = {
        "\\bi8\\b", "\\bi16\\b", "\\bi32\\b", "\\bi64\\b", "\\bi128\\b",
        "\\bu8\\b", "\\bu16\\b", "\\bu32\\b", "\\bu64\\b", "\\bu128\\b",
        "\\bfloat\\b", "\\bf32\\b", "\\bf64\\b", "\\bbool\\b", "\\bchar\\b",
        "\\bstr\\b", "\\busize\\b", "\\bisize\\b"
    };
    
    for (const QString &pattern : typePatterns) {
        HighlightingRule rule;
        rule.pattern = QRegularExpression(pattern);
        rule.format = typeFormat_;
        highlightingRules_.append(rule);
    }
    
    HighlightingRule rule;
    
    rule.pattern = QRegularExpression("\"([^\"\\\\]|\\\\.)*\"");
    rule.format = stringFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("\\b[0-9]+\\b");
    rule.format = numberFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("//[^\n]*");
    rule.format = commentFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("#\\[\\w+.*\\]");
    rule.format = preprocessorFormat_;
    highlightingRules_.append(rule);
    
    commentStartExpression_ = QRegularExpression("/\\*");
    commentEndExpression_ = QRegularExpression("\\*/");
}

void SyntaxHighlighter::setupMakefileLanguage()
{
    HighlightingRule rule;
    
    rule.pattern = QRegularExpression("^[A-Za-z_][A-Za-z0-9_]*:");
    rule.format = functionFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("\\$\\([A-Za-z_][A-Za-z0-9_]*\\)");
    rule.format = typeFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("#[^\n]*");
    rule.format = commentFormat_;
    highlightingRules_.append(rule);
}

void SyntaxHighlighter::setupShellLanguage()
{
    QStringList keywordPatterns = {
        "\\bif\\b", "\\bthen\\b", "\\belse\\b", "\\belif\\b", "\\bfi\\b",
        "\\bfor\\b", "\\bwhile\\b", "\\bdo\\b", "\\bdone\\b", "\\bcase\\b",
        "\\besac\\b", "\\bfunction\\b", "\\breturn\\b", "\\bexit\\b"
    };
    
    for (const QString &pattern : keywordPatterns) {
        HighlightingRule rule;
        rule.pattern = QRegularExpression(pattern);
        rule.format = keywordFormat_;
        highlightingRules_.append(rule);
    }
    
    HighlightingRule rule;
    
    rule.pattern = QRegularExpression("\"([^\"\\\\]|\\\\.)*\"");
    rule.format = stringFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("'([^'\\\\]|\\\\.)*'");
    rule.format = stringFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("\\$[A-Za-z_][A-Za-z0-9_]*");
    rule.format = typeFormat_;
    highlightingRules_.append(rule);
    
    rule.pattern = QRegularExpression("#[^\n]*");
    rule.format = commentFormat_;
    highlightingRules_.append(rule);
}

void SyntaxHighlighter::highlightBlock(const QString &text)
{
    for (const HighlightingRule &rule : std::as_const(highlightingRules_)) {
        QRegularExpressionMatchIterator matchIterator = rule.pattern.globalMatch(text);
        while (matchIterator.hasNext()) {
            QRegularExpressionMatch match = matchIterator.next();
            setFormat(match.capturedStart(), match.capturedLength(), rule.format);
        }
    }
    
    setCurrentBlockState(0);
    
    if (language_ == C || language_ == CPlusPlus || language_ == Rust) {
        int startIndex = 0;
        if (previousBlockState() != 1)
            startIndex = text.indexOf(commentStartExpression_);
        
        while (startIndex >= 0) {
            QRegularExpressionMatch match = commentEndExpression_.match(text, startIndex);
            int endIndex = match.capturedStart();
            int commentLength = 0;
            if (endIndex == -1) {
                setCurrentBlockState(1);
                commentLength = text.length() - startIndex;
            } else {
                commentLength = endIndex - startIndex + match.capturedLength();
            }
            setFormat(startIndex, commentLength, commentFormat_);
            startIndex = text.indexOf(commentStartExpression_, startIndex + commentLength);
        }
    }
}
