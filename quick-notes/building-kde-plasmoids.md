# Getting started building basic KDE plasmoids

## What is a plasmoid?

'Plasmoid' is a fancy name to describe a KDE Plasma widget.
System tray, system clock, launch menu, attached to either the desktop or a taskbar... All of these are plasmoids shipped with Plasma!   

But what if you need to write one yourself? How does it work?  
I recently decided to write one and was stuck on the [official get started documentation](https://develop.kde.org/docs/plasma/widget/setup/) for quite a while. There is a lot of implicit behavior and things are not as obvious as I would like. I probably won't need to write another one for a while, and by that time, it's very likely I will need to start over entirely...
Therefore, I decided to write my own guide, including information I did not see in the documentation, to make sure I don't get stuck as long the next time around!

## File structure

The file structure described in the documentation is not a suggestion. There are some implicit expectations for this file structure.

```txt
└── ~/Code/plasmoid-helloworld/
    └── package
        ├── contents
        │   ├── ui
        │   │   ├── main.qml
        │   │   ├── CompactRepresentation.qml
        │   │   ├── FullRepresentation.qml
        │   │   └── configGeneral.qml
        │   └── config
        │       ├── config.qml
        │       └── main.xml
        └── metadata.json
```

## Metadata

Every plasmoid requires a `metadata.json` file.
The `metadata.desktop` that can still be found in many apps and plasmoids has been deprecated, and the JSON version should now be preferred.

```json
{
	"KPlugin": {
		"Authors": [
			{
				"Email": "contact@ade-sede.dev",
				"Name": "Adrien DE SEDE"
			}
		],
		"Category": "Date and Time",
		"Description": "Displays transit times for configured routes",
		"Icon": "office-calendar",
		"Id": "com.github.ade-sede.home-dashboard.transit-clock",
		"Name": "Transit Clock",
		"Version": "1.0"
	},
	"X-Plasma-API": "declarativeappletscript",
	"X-Plasma-MainScript": "ui/main.qml",
	"KPackageStructure": "Plasma/Applet"
}
```

Documentation of all properties can be found in the [developer documentation](https://develop.kde.org/docs/plasma/widget/properties/#metadatajson).
Notably:
- `KPackageStructure` is what dictates the 'implicit expectations' on the file structure
- Notice how `X-Plasma-MainScript` is `ui/main.qml` and not `contents/ui/main.qml`. Good luck getting your plasmoid to work if you don't know that!

## QML

QML is a [simple declarative language](https://doc.qt.io/qt-6/qmlapplications.html) used to describe the layout that is then rendered by Qt.  
The file `main.qml` is typically the 'root' of the plasmoid. This is where you would define common building blocks:
- variables
- signals (used for various 'containers' and plasmoids to communicate)
- functions (more or less written in JavaScript)

But also where you would define the full and the compact representation of your plasmoid.  
The _compact_ representation is the one used when the plasmoid is attached to a taskbar.  
The _full_ representation is the one used when you click on the taskbar or if you pin the plasmoid to a desktop.  

```qml
Plasmoid.compactRepresentation: CompactRepresentation {}
Plasmoid.fullRepresentation: FullRepresentation {}
```

In your `main.qml` file, these two lines would mean that the file `CompactRepresentation.qml` is what is used for the _compact_ representation and `FullRepresentation.qml` for the _full_ representation. These files can access common blocks defined in `main.qml` through `plasmoid.rootItem`, making it easy to share data and functions across your widget's different views.

## Views and Containers in QML

QML uses a hierarchical system of views and containers to organize UI elements:

### Containers
Containers are QML elements that can hold and arrange other elements. Common containers include:

- `Item`: The basic invisible container, used as a generic holder for other components
- `Rectangle`: A visible container with properties for color, border, etc.
- `Column`, `Row`, and `Grid`: Layout containers that arrange their children in specific patterns
- `ColumnLayout`, `RowLayout`, and `GridLayout`: More powerful layouts with additional alignment and sizing options

Containers help organize your UI both visually and in the code structure. For example:

```qml
Rectangle {
    id: mainContainer
    width: 200
    height: 200
    color: "lightgray"
    
    Column {
        spacing: 10
        anchors.centerIn: parent
        
        Rectangle { width: 100; height: 30; color: "red" }
        Rectangle { width: 100; height: 30; color: "blue" }
    }
}
```

### Views
Views are special containers designed to display collections of data:

- `ListView`: Displays a list of items based on a data model
- `GridView`: Arranges items in a grid pattern
- `PathView`: Arranges items along a path
- `Repeater`: A simpler component that creates multiple instances of an item

Views in plasmoids are particularly useful for presenting dynamic data like notifications, system information, or in my case, transit schedules.

The connection between containers, views, and the plasmoid structure is important - your `main.qml` serves as the root container, while your representation files define specific views of your widget depending on context.

## Config

Your plasmoid may need to store some state related to configuration.
You can define 'fields' that are part of this configuration in `config/main.xml`.

```xml
<?xml version="1.0" encoding="UTF-8"?>
<kcfg xmlns="http://www.kde.org/standards/kcfg/1.0"
      xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
      xsi:schemaLocation="http://www.kde.org/standards/kcfg/1.0
      http://www.kde.org/standards/kcfg/1.0/kcfg.xsd" >
  <kcfgfile name=""/>

  <group name="General">
    <entry name="serverUrl" type="String">
      <default>https://some-domain-name.com</default>
    </entry>
    <entry name="username" type="String">
      <default>some-username</default>
    </entry>
  </group>
</kcfg>
```

You can access them from your QML files using the following syntax:
```qml
Item {
    property string serverUrl: plasmoid.configuration.serverUrl
    property string username: plasmoid.configuration.username
}
```

You can then define what the config menu looks like using `config/config.qml`.

```qml
import QtQuick 2.0
import org.kde.plasma.configuration 2.0

ConfigModel {
    ConfigCategory {
        name: i18n("General")
        icon: "configure"
        source: "ConfigGeneral.qml"
    }
}
```

Each `ConfigCategory` defines a 'pane' of the configuration menu.  
The property `source` links to the file that is used to display the 'content' of that pane. Again, note that it is implicitly assumed the file is actually `contents/ui/ConfigGeneral.qml`.

```qml
Item {
    property alias cfg_serverUrl: serverUrlField.text
}
```

Setting that property 'auto-magically' updates the `serverUrl` property that we defined earlier. I have no idea where this state is stored... but the Plasma framework handles this persistence for you.

The configuration data is actually stored in `~/.config/plasma-org.kde.plasma.desktop-appletsrc`, with your plasmoid's settings in a section named after your plasmoid's ID. For example:

```
[Containments][XX][Applets][YY][Configuration]
serverUrl=https://some-domain-name.com
username=some-username
```

Where XX and YY are numeric IDs assigned by Plasma. The KDE framework handles all the reading and writing to this file through KConfig and KConfigGroup, which is why you don't need to manage it yourself.

## Installing your plasmoid

The entirety of the `package` directory must be copied to `~/.local/share/plasma/plasmoids/<id>` where id matches what we set in the `metadata.json` file.  
In our case, it would be `~/.local/share/plasma/plasmoids/com.github.ade-sede.home-dashboard.transit-clock`.
The simplest way to get it installed is to run `plasmapkg2 -t Plasma/Applet -i ./package`.

For the plasmoid to be available in the widget menu, you probably need to log out and log back in.

## Testing your plasmoid

`plasmoidviewer -a ./package/` will launch a miniature version of the plasma shell where your plasmoid is open in its full representation.

To test the compact representation (as it would appear in a taskbar), use:
```
plasmoidviewer -f horizontal -a ./package/
```

The `-f` (--formfactor) parameter with "horizontal" simulates a panel environment where the compact representation would be used.

You can also create a test panel with:
```
plasmoidviewer -c org.kde.panel -a ./package/
```

This creates a panel container where your plasmoid will default to its compact representation.

## Learning QML syntax, patterns, etc.

Reading through some parts of the [plasma-applet-eventcalendar](https://github.com/Zren/plasma-applet-eventcalendar) was useful.
I've also found asking LLMs for basic examples to be quite effective.
