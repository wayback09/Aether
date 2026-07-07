# UI Layout and Components

## Sidebar
The sidebar is permanent.
It never collapses automatically.
Width: 240px

Contains:
- Logo
- Navigation
- Extensions
- Settings
- Account

Example:
--------------------------------
LOGO

Home
Instances
Extensions
-------------
(Settings icon)
(User)
--------------------------------

Extensions appear dynamically.
The launcher itself never contains built-in feature pages.

## Home
Purpose: Launch Minecraft. Nothing else.

Contains:
- Current Instance
- Play Button
- Minecraft Version
- Loader
- RAM Allocation
- Recent Activity
- Account

No news.
No changelog.
No advertisements.
No promotions.
No banners.
No Discord widgets.

Layout:
---------------------------------
Home

Current Instance
Play
Version
Loader
Memory
---------------------------------
Large whitespace below.

## Instances
Purpose: Manage Minecraft instances.

Actions:
- Create
- Duplicate
- Rename
- Delete
- Import
- Export

Each instance card shows:
- Name
- Minecraft Version
- Loader
- Last Played
- Play Button

No statistics.
No graphs.
No unnecessary metadata.

Instance Details Tabs:
- Play
- Settings
- Extensions
Nothing else.
Mods, Logs, Servers, Worlds are added only by extensions.

## Extensions
Purpose: Manage installed extensions.

Contains:
- Installed Extensions
- Updates
- Permissions
- Enable
- Disable
- Remove
- Restart Extension
- Resource Usage
- Health

Every extension card shows:
- Icon
- Name
- Version
- Author
- Status
- Memory Usage
- CPU Usage
- Permissions
- Restart Button

## Settings
Categories:
- Launcher
- Appearance
- Java
- Updates
- Advanced
- Extensions

Each category is simple. Avoid overwhelming users.

## Non-Existent Built-in Pages
The following pages DO NOT EXIST in the core launcher and are created by extensions:
- Downloads (Created by the official Minecraft extension)
- Modrinth (Created by Modrinth extension)
- CurseForge (Created by CurseForge extension)
- Logs (Created by Log Viewer extension)
- Server Browser (Created by Server Browser extension)
- Worlds (Created by World Manager extension)
- Skin Manager (Created by Skin extension)

## Empty States
Every page should have beautiful empty states.

Example:
No Extensions Installed
Install your first extension to add new functionality.
[Browse Extensions]

Never leave blank pages.

## Dialogs
Centered. Simple. One action. One cancel.
Never overload dialogs.

## Notifications
Bottom Right. Disappear automatically. Used sparingly.

## Progress
Downloads, Launching, Installing should always show progress.
Never leave users guessing.

## Responsive Behaviour
Minimum Width: 1100px.
Never become mobile. Desktop first.

## UI Consistency Rules
Every page title appears in the same location.
Every toolbar behaves the same.
Every search bar behaves the same.
Every table behaves the same.
Every settings page behaves the same.
Every dialog behaves the same.
No extension may violate these rules.
