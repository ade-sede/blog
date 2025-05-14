# Building KDE Plasmoids for Plasma 6: Notes

_Mostly for myself as a form of internal documentation..._

## What is a plasmoid?

'Plasmoid' is a fancy name for a KDE Plasma widget. System tray, clock, application launcher - these are all plasmoids shipped with Plasma.

I recently created one and hit several roadblocks with the [official documentation](https://develop.kde.org/docs/plasma/widget/setup/). Here are my notes on the implicit behaviors and non-obvious aspects I discovered.

## File structure

Plasma has strict expectations about file structure. Deviating from this pattern leads to mysterious failures:

```directory-structure
package/
  metadata.json
  contents/
    ui/
      main.qml
      CompactRepresentation.qml
      FullRepresentation.qml
      configGeneral.qml
    config/
      config.qml
      main.xml
```

## Metadata

The `metadata.json` file defines the plasmoid's identity and behavior:

```json:package/metadata.json
{
    "KPlugin": {
        "Authors": [
            {
                "Email": "contact@example.com",
                "Name": "Your Name"
            }
        ],
        "Category": "Utilities",
        "Description": "A simple Hello World plasmoid",
        "Icon": "document-text",
        "Id": "com.example.helloworld",
        "Name": "Hello World",
        "Version": "1.0"
    },
    "X-Plasma-API": "declarativeappletscript",
    "X-Plasma-MainScript": "ui/main.qml",
    "KPackageStructure": "Plasma/Applet"
}
```

Things that tripped me up:
- `KPackageStructure` dictates the expected file structure
- `X-Plasma-MainScript` is relative to the `contents/` directory, not the `package/` directory 
- The `Id` becomes the installation directory name under `~/.local/share/plasma/plasmoids/`

## QML Basics

QML powers the UI of plasmoids. The main.qml file serves as the entry point:

```qml:package/contents/ui/main.qml
import QtQuick 2.0
import org.kde.plasma.plasmoid 2.0
import org.kde.plasma.core 2.0 as PlasmaCore

Item {
    id: root
    
    property string greeting: "Hello, Plasma!"
    
    Plasmoid.fullRepresentation: Rectangle {
        width: 200
        height: 100
        color: PlasmaCore.Theme.backgroundColor
        
        Text {
            anchors.centerIn: parent
            text: root.greeting
            color: PlasmaCore.Theme.textColor
        }
    }
}
```

Using `PlasmaCore.Theme` properties automatically adapts your plasmoid to the user's color scheme, which is pretty handy.

## Compact and Full Representations

Plasmoids have two display modes: compact (for panels) and full (for expanded views or desktop widgets).

For panels, you need a compact representation:

```qml:package/contents/ui/CompactRepresentation.qml
import QtQuick 2.0
import org.kde.plasma.core 2.0 as PlasmaCore

Item {
    id: compactRoot
    
    width: 24
    height: 24
    
    Rectangle {
        anchors.fill: parent
        color: PlasmaCore.Theme.highlightColor
        
        Text {
            anchors.centerIn: parent
            text: "HW"
            color: PlasmaCore.Theme.highlightedTextColor
            font.pixelSize: parent.height * 0.6
        }
    }
    
    MouseArea {
        anchors.fill: parent
        onClicked: plasmoid.expanded = !plasmoid.expanded
    }
}
```

For the expanded view, you use a full representation:

```qml:package/contents/ui/FullRepresentation.qml
import QtQuick 2.0
import org.kde.plasma.core 2.0 as PlasmaCore

Rectangle {
    id: fullRoot
    
    width: 200
    height: 100
    color: PlasmaCore.Theme.backgroundColor
    
    Text {
        anchors.centerIn: parent
        text: plasmoid.rootItem.greeting
        color: PlasmaCore.Theme.textColor
    }
}
```

When splitting these into separate files, the main.qml references them:

```qml:package/contents/ui/main.qml
import QtQuick 2.0
import org.kde.plasma.plasmoid 2.0
import org.kde.plasma.core 2.0 as PlasmaCore

Item {
    id: root
    
    property string greeting: "Hello, Plasma!"
    
    Plasmoid.compactRepresentation: CompactRepresentation {}
    Plasmoid.fullRepresentation: FullRepresentation {}
}
```

I discovered that these files can access data from main.qml through `plasmoid.rootItem` - useful for sharing variables between different views.

## Configuration System

The configuration system involves three files working together:

The `main.xml` file defines what settings you want to save:

```xml:package/contents/config/main.xml
<?xml version="1.0" encoding="UTF-8"?>
<kcfg xmlns="http://www.kde.org/standards/kcfg/1.0"
      xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
      xsi:schemaLocation="http://www.kde.org/standards/kcfg/1.0
      http://www.kde.org/standards/kcfg/1.0/kcfg.xsd" >
  <kcfgfile name=""/>

  <group name="General">
    <entry name="customGreeting" type="String">
      <default>Hello, Plasma!</default>
    </entry>
  </group>
</kcfg>
```

The `config.qml` file defines how the settings UI is organized:

```qml:package/contents/config/config.qml
import QtQuick 2.0
import org.kde.plasma.configuration 2.0

ConfigModel {
    ConfigCategory {
        name: "General"
        icon: "configure"
        source: "configGeneral.qml"
    }
}
```

Something weird: `source` points to `contents/ui/configGeneral.qml`, not to `contents/config/configGeneral.qml` as you might expect.

The `configGeneral.qml` file is the actual settings form:

```qml:package/contents/ui/configGeneral.qml
import QtQuick 2.0
import QtQuick.Controls 2.5
import QtQuick.Layouts 1.12

Item {
    id: configPage
    
    property alias cfg_customGreeting: greetingField.text
    
    ColumnLayout {
        anchors.left: parent.left
        anchors.right: parent.right
        
        Label {
            text: "Custom greeting:"
        }
        
        TextField {
            id: greetingField
            Layout.fillWidth: true
        }
    }
}
```

The naming convention `cfg_` + property name is how KDE knows which properties to save. When you edit the text field, it automatically saves the value - no extra code needed.

In `main.qml`, you access saved settings with:

```qml
property string greeting: plasmoid.configuration.customGreeting
```

The values are stored in `~/.config/plasma-org.kde.plasma.desktop-appletsrc`, but you never need to touch this file directly.

## Testing and Installation

The `plasmoidviewer` tool is great for testing:

```fish
# View full representation
plasmoidviewer -a ~/Code/plasmoid-helloworld/package/

# View compact representation (as in panel)
plasmoidviewer -f horizontal -a ~/Code/plasmoid-helloworld/package/

# Test in a panel container
plasmoidviewer -c org.kde.panel -a ~/Code/plasmoid-helloworld/package/
```

To install, you copy the package directory to `~/.local/share/plasma/plasmoids/<id>`:

```fish
# First-time installation
plasmapkg2 -t Plasma/Applet -i ~/Code/plasmoid-helloworld/package/

# Update existing installation
plasmapkg2 -t Plasma/Applet -u ~/Code/plasmoid-helloworld/package/
```

You'll need to log out and back in before the plasmoid shows up in the widget selection menu.

## QML Components

QML has different types of layout elements:

**Containers** organize and arrange elements:
- `Item`: Invisible container (the basic one)
- `Rectangle`: Visible container with color, border, etc.
- `Column`, `Row`, and `Grid`: Simple layout containers
- `ColumnLayout`, `RowLayout`, and `GridLayout`: More powerful layout containers

**Views** display collections of data:
- `ListView`: Shows items in a list
- `GridView`: Shows items in a grid
- `PathView`: Shows items along a path
- `Repeater`: Simple way to create multiple copies of an item

## Learning More

I found the code from [plasma-applet-eventcalendar](https://github.com/Zren/plasma-applet-eventcalendar) more helpful than the official docs in many cases.

I've also had good luck asking LLMs for QML examples when I get stuck.

## But ... Plasma 6

A few weeks after writing this quicknote I upgraded to Plasma 6 and you guessed it: my widget does not work...
I thought I was smart taking notes for later in case I needed to do it again, turns out I will need to re-learn from scratch anyway...
