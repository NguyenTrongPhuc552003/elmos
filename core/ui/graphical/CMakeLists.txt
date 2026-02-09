cmake_minimum_required(VERSION 3.16)

project(elmos-gui VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)

# Find Qt6 (adjust path for Homebrew on macOS)
set(CMAKE_PREFIX_PATH "/opt/homebrew/opt/qt")
find_package(Qt6 REQUIRED COMPONENTS Core Widgets Gui Network)

# gRPC (includes Protobuf transitively)
find_package(gRPC CONFIG REQUIRED)

# Generated proto files (pre-generated with protoc)
set(GENERATED_PATH "${CMAKE_SOURCE_DIR}/generated/api/proto")
set(PROTO_GENERATED_SOURCES
    ${GENERATED_PATH}/common.pb.cc
    ${GENERATED_PATH}/common.grpc.pb.cc
    ${GENERATED_PATH}/kernel.pb.cc
    ${GENERATED_PATH}/kernel.grpc.pb.cc
    ${GENERATED_PATH}/workspace.pb.cc
    ${GENERATED_PATH}/workspace.grpc.pb.cc
    ${GENERATED_PATH}/toolchain.pb.cc
    ${GENERATED_PATH}/toolchain.grpc.pb.cc
    ${GENERATED_PATH}/qemu.pb.cc
    ${GENERATED_PATH}/qemu.grpc.pb.cc
)

# Source files
set(SOURCES
    src/main.cpp
    src/mainwindow.cpp
    src/codeeditor.cpp
    src/syntaxhighlighter.cpp
    src/projectwizard.cpp
    src/projectexplorer.cpp
    src/grpcclient.cpp
    src/kernelbuildwidget.cpp
    src/qemuconsolewidget.cpp
    src/settingsdialog.cpp
    src/toolchainwidget.cpp
    src/workspacewidget.cpp
    src/archselectorwidget.cpp
    ${PROTO_GENERATED_SOURCES}
)

set(HEADERS
    include/mainwindow.h
    include/codeeditor.h
    include/syntaxhighlighter.h
    include/projectwizard.h
    include/projectexplorer.h
    include/grpcclient.h
    include/kernelbuildwidget.h
    include/qemuconsolewidget.h
    include/settingsdialog.h
    include/toolchainwidget.h
    include/workspacewidget.h
    include/archselectorwidget.h
)

# UI forms (to be added later)
# set(FORMS
#     forms/mainwindow.ui
#     forms/projectwizard.ui
# )

# Resources (to be added later)
# set(RESOURCES
#     resources/resources.qrc
# )

# Include directories
include_directories(
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${CMAKE_CURRENT_SOURCE_DIR}/generated
    ${CMAKE_CURRENT_BINARY_DIR}
)

# Executable
add_executable(elmos-gui
    ${SOURCES}
    ${HEADERS}
)

# Link libraries
target_link_libraries(elmos-gui
    Qt6::Core
    Qt6::Widgets
    Qt6::Gui
    Qt6::Network
    gRPC::grpc++
    gRPC::grpc++_reflection
)

# Install target
install(TARGETS elmos-gui
    RUNTIME DESTINATION bin
    BUNDLE DESTINATION .
)

# macOS bundle configuration
if(APPLE)
    set_target_properties(elmos-gui PROPERTIES
        MACOSX_BUNDLE TRUE
        MACOSX_BUNDLE_GUI_IDENTIFIER "com.elmos.gui"
        MACOSX_BUNDLE_BUNDLE_NAME "ELMOS"
        MACOSX_BUNDLE_BUNDLE_VERSION ${PROJECT_VERSION}
        MACOSX_BUNDLE_SHORT_VERSION_STRING ${PROJECT_VERSION}
    )
endif()

# Windows configuration
if(WIN32)
    set_target_properties(elmos-gui PROPERTIES
        WIN32_EXECUTABLE TRUE
    )
endif()
