#include "qemuconsolewidget.h"
#include "grpcclient.h"
#include <QGroupBox>
#include <QScrollBar>
#include <QThread>
#include <QRegularExpression>
#include <QKeyEvent>

QEMUConsoleWidget::QEMUConsoleWidget(GrpcClient *grpcClient, QWidget *parent)
    : QWidget(parent)
    , grpcClient_(grpcClient)
    , isRunning_(false)
{
    setupUI();
    connectSignals();

    defaultFormat_.setForeground(Qt::white);
    defaultFormat_.setBackground(Qt::black);
    currentFormat_ = defaultFormat_;
}

QEMUConsoleWidget::~QEMUConsoleWidget() = default;

void QEMUConsoleWidget::setupUI()
{
    QVBoxLayout *mainLayout = new QVBoxLayout(this);

    QGroupBox *configGroup = new QGroupBox("QEMU Configuration");
    QHBoxLayout *configLayout = new QHBoxLayout(configGroup);

    graphicalCheckBox_ = new QCheckBox("Graphical");
    configLayout->addWidget(graphicalCheckBox_);

    debugCheckBox_ = new QCheckBox("Debug (GDB:1234)");
    configLayout->addWidget(debugCheckBox_);

    configLayout->addWidget(new QLabel("Memory (MB):"));
    memorySpinBox_ = new QSpinBox();
    memorySpinBox_->setRange(0, 16384);
    memorySpinBox_->setValue(256);
    memorySpinBox_->setSpecialValueText("Auto");
    configLayout->addWidget(memorySpinBox_);

    configLayout->addWidget(new QLabel("CPUs:"));
    cpusSpinBox_ = new QSpinBox();
    cpusSpinBox_->setRange(0, 16);
    cpusSpinBox_->setValue(2);
    cpusSpinBox_->setSpecialValueText("Auto");
    configLayout->addWidget(cpusSpinBox_);

    configLayout->addStretch();

    mainLayout->addWidget(configGroup);

    QHBoxLayout *cmdlineLayout = new QHBoxLayout();
    cmdlineLayout->addWidget(new QLabel("Kernel cmdline:"));
    cmdlineEdit_ = new QLineEdit();
    cmdlineEdit_->setPlaceholderText("console=ttyAMA0");
    cmdlineLayout->addWidget(cmdlineEdit_);
    mainLayout->addLayout(cmdlineLayout);

    consoleView_ = new QTextEdit();
    consoleView_->setReadOnly(true);
    consoleView_->setFont(QFont("Monaco", 11));
    consoleView_->setStyleSheet("background-color: #000; color: #fff;");
    consoleView_->setLineWrapMode(QTextEdit::WidgetWidth);
    mainLayout->addWidget(consoleView_, 1);

    inputLine_ = new QLineEdit();
    inputLine_->setEnabled(false);
    inputLine_->setPlaceholderText("Type commands here (press Enter to send)");
    inputLine_->setFont(QFont("Monaco", 11));
    connect(inputLine_, &QLineEdit::returnPressed, this, &QEMUConsoleWidget::onReturnPressed);
    mainLayout->addWidget(inputLine_);

    QHBoxLayout *buttonLayout = new QHBoxLayout();
    startButton_ = new QPushButton("Start QEMU");
    stopButton_ = new QPushButton("Stop");
    clearButton_ = new QPushButton("Clear");
    stopButton_->setEnabled(false);

    buttonLayout->addWidget(startButton_);
    buttonLayout->addWidget(stopButton_);
    buttonLayout->addWidget(clearButton_);
    buttonLayout->addStretch();

    statusLabel_ = new QLabel("Ready");
    buttonLayout->addWidget(statusLabel_);

    mainLayout->addLayout(buttonLayout);

    connect(startButton_, &QPushButton::clicked, this, &QEMUConsoleWidget::startQEMU);
    connect(stopButton_, &QPushButton::clicked, this, &QEMUConsoleWidget::stopQEMU);
    connect(clearButton_, &QPushButton::clicked, this, &QEMUConsoleWidget::clearConsole);
}

void QEMUConsoleWidget::connectSignals()
{
    if (!grpcClient_) {
        return;
    }

    connect(grpcClient_, &GrpcClient::qemuStarted, this, &QEMUConsoleWidget::onQEMUStarted);
    connect(grpcClient_, &GrpcClient::qemuConsoleOutput, this, &QEMUConsoleWidget::onQEMUConsoleOutput);
    connect(grpcClient_, &GrpcClient::qemuStopped, this, &QEMUConsoleWidget::onQEMUStopped);
    connect(grpcClient_, &GrpcClient::qemuError, this, &QEMUConsoleWidget::onQEMUError);
}

void QEMUConsoleWidget::startQEMU()
{
    if (isRunning_ || !grpcClient_) {
        return;
    }

    isRunning_ = true;
    startButton_->setEnabled(false);
    stopButton_->setEnabled(true);
    inputLine_->setEnabled(true);
    inputLine_->setFocus();
    statusLabel_->setText("Starting...");
    consoleView_->clear();

    bool graphical = graphicalCheckBox_->isChecked();
    bool debug = debugCheckBox_->isChecked();
    int memory = memorySpinBox_->value();
    int cpus = cpusSpinBox_->value();
    QString cmdline = cmdlineEdit_->text();

    QThread *qemuThread = QThread::create([this, graphical, debug, memory, cpus, cmdline]() {
        grpcClient_->runQEMU(graphical, debug, memory, cpus, QStringList(), cmdline);
    });
    
    connect(qemuThread, &QThread::finished, qemuThread, &QThread::deleteLater);
    qemuThread->start();
}

void QEMUConsoleWidget::stopQEMU()
{
    if (!isRunning_ || !grpcClient_) {
        return;
    }

    grpcClient_->stopQEMU();
    isRunning_ = false;
    startButton_->setEnabled(true);
    stopButton_->setEnabled(false);
    inputLine_->setEnabled(false);
    statusLabel_->setText("Stopped");
}

void QEMUConsoleWidget::clearConsole()
{
    consoleView_->clear();
}

void QEMUConsoleWidget::sendInput()
{
    if (!isRunning_ || !grpcClient_) {
        return;
    }

    QString text = inputLine_->text();
    if (text.isEmpty()) {
        return;
    }

    QByteArray data = (text + "\n").toUtf8();
    grpcClient_->sendQEMUInput(data);
    
    appendText(text + "\n", currentFormat_);
    inputLine_->clear();
}

void QEMUConsoleWidget::onQEMUStarted(int pid, const QString &qemuVersion, const QString &command)
{
    statusLabel_->setText(QString("Running (PID: %1)").arg(pid));
    
    QTextCharFormat greenFormat;
    greenFormat.setForeground(Qt::green);
    greenFormat.setBackground(Qt::black);
    
    appendText(QString("=== QEMU Started ===\n"), greenFormat);
    appendText(QString("Version: %1\n").arg(qemuVersion), greenFormat);
    appendText(QString("PID: %1\n\n").arg(pid), greenFormat);
}

void QEMUConsoleWidget::onQEMUConsoleOutput(const QByteArray &data)
{
    processANSI(data);
}

void QEMUConsoleWidget::onQEMUStopped(int exitCode, qint64 uptimeMs)
{
    isRunning_ = false;
    startButton_->setEnabled(true);
    stopButton_->setEnabled(false);
    inputLine_->setEnabled(false);
    
    statusLabel_->setText(QString("Stopped (exit: %1)").arg(exitCode));
    
    QTextCharFormat yellowFormat;
    yellowFormat.setForeground(Qt::yellow);
    yellowFormat.setBackground(Qt::black);
    
    appendText(QString("\n=== QEMU Stopped ===\n"), yellowFormat);
    appendText(QString("Exit code: %1\n").arg(exitCode), yellowFormat);
    appendText(QString("Uptime: %1s\n").arg(uptimeMs / 1000.0, 0, 'f', 2), yellowFormat);
}

void QEMUConsoleWidget::onQEMUError(const QString &error)
{
    QTextCharFormat redFormat;
    redFormat.setForeground(Qt::red);
    redFormat.setBackground(Qt::black);
    
    appendText(QString("[ERROR] %1\n").arg(error), redFormat);
}

void QEMUConsoleWidget::onReturnPressed()
{
    sendInput();
}

void QEMUConsoleWidget::processANSI(const QByteArray &data)
{
    static QRegularExpression ansiRegex("\\x1b\\[([0-9;]*)([a-zA-Z])");
    
    QString text = QString::fromUtf8(data);
    int pos = 0;
    
    while (pos < text.length()) {
        QRegularExpressionMatch match = ansiRegex.match(text, pos);
        
        if (match.hasMatch() && match.capturedStart() == pos) {
            QString params = match.captured(1);
            QChar command = match.captured(2)[0];
            
            if (command == 'm') {
                QStringList codes = params.split(';', Qt::SkipEmptyParts);
                for (const QString &code : codes) {
                    int c = code.toInt();
                    
                    if (c == 0) {
                        currentFormat_ = defaultFormat_;
                    } else if (c >= 30 && c <= 37) {
                        QColor colors[] = {Qt::black, Qt::red, Qt::green, Qt::yellow, 
                                          Qt::blue, Qt::magenta, Qt::cyan, Qt::white};
                        currentFormat_.setForeground(colors[c - 30]);
                    } else if (c >= 40 && c <= 47) {
                        QColor colors[] = {Qt::black, Qt::red, Qt::green, Qt::yellow, 
                                          Qt::blue, Qt::magenta, Qt::cyan, Qt::white};
                        currentFormat_.setBackground(colors[c - 40]);
                    } else if (c == 1) {
                        currentFormat_.setFontWeight(QFont::Bold);
                    }
                }
            }
            
            pos = match.capturedEnd();
        } else {
            int nextEscape = text.indexOf('\x1b', pos);
            if (nextEscape == -1) {
                appendText(text.mid(pos), currentFormat_);
                break;
            } else {
                appendText(text.mid(pos, nextEscape - pos), currentFormat_);
                pos = nextEscape;
            }
        }
    }
}

void QEMUConsoleWidget::appendText(const QString &text, const QTextCharFormat &format)
{
    QTextCursor cursor = consoleView_->textCursor();
    cursor.movePosition(QTextCursor::End);
    cursor.setCharFormat(format);
    cursor.insertText(text);
    
    QScrollBar *scrollBar = consoleView_->verticalScrollBar();
    scrollBar->setValue(scrollBar->maximum());
}
