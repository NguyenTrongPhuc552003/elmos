#include "mainwindow.h"
#include <QApplication>
#include <QStyleFactory>
#include <QFont>

int main(int argc, char *argv[])
{
    QApplication app(argc, argv);
    
    // Application metadata
    app.setOrganizationName("ELMOS");
    app.setOrganizationDomain("elmos.dev");
    app.setApplicationName("ELMOS IDE");
    app.setApplicationVersion("1.0.0");
    
    // Set modern dark theme (VSCode-like)
    app.setStyle(QStyleFactory::create("Fusion"));
    
    QPalette darkPalette;
    darkPalette.setColor(QPalette::Window, QColor(30, 30, 30));
    darkPalette.setColor(QPalette::WindowText, Qt::white);
    darkPalette.setColor(QPalette::Base, QColor(25, 25, 25));
    darkPalette.setColor(QPalette::AlternateBase, QColor(45, 45, 45));
    darkPalette.setColor(QPalette::ToolTipBase, Qt::white);
    darkPalette.setColor(QPalette::ToolTipText, Qt::white);
    darkPalette.setColor(QPalette::Text, Qt::white);
    darkPalette.setColor(QPalette::Button, QColor(45, 45, 45));
    darkPalette.setColor(QPalette::ButtonText, Qt::white);
    darkPalette.setColor(QPalette::BrightText, Qt::red);
    darkPalette.setColor(QPalette::Link, QColor(42, 130, 218));
    darkPalette.setColor(QPalette::Highlight, QColor(42, 130, 218));
    darkPalette.setColor(QPalette::HighlightedText, Qt::black);
    app.setPalette(darkPalette);
    
    // Set monospace font for consistency (using macOS compatible fonts)
    QFont font;
    font.setFamily("Monaco");  // macOS default monospace font
    font.setPointSize(11);
    app.setFont(font);
    
    // Create and show main window
    MainWindow mainWindow;
    mainWindow.setWindowTitle("ELMOS - Embedded Linux Development IDE");
    mainWindow.setWindowIcon(QIcon(":/icons/elmos-logo.png"));
    mainWindow.resize(1400, 900);
    mainWindow.show();
    
    return app.exec();
}
