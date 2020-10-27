Very brief instructions. See the Meerkat-API section on our forums
for help (http://poker-genius.com/forum/Thread-Meerkat-API-help)

----------------------------------------------------------------------
Requires Poker Genius v. 1.05 or higher
----------------------------------------------------------------------

Steps to build a bot:

* Implement the poker.Player interface, and over-ride the required methods.

* Build against the provided meerkat-api.jar library

* Create a jar file with your bot's classes.

* Create a player definition file (see example)

* Drop both the .pd and .jar file into the Program's data/bots/ folder
  The default location on windows is C:\Program Files\PokerGenius\data\bots 
  or C:\Program Files (x86)\PokerGenius\data\bots
  On Mac OS X you must install it inside the application bundle, which is
  in Contents\Resources\Java\data\bots\

* Load the Opponent Manager and import a new Opponent (your new bot type
  should appear in the list of installed AI Engines)

* See the SimpleBot.java example for a basic bot.

* # All plug-ins must run through the following class:
PLAYER_CLASS=PlugInOpponent


----------------------------------------------------------------------
 Specifying a No-Limit Bot
----------------------------------------------------------------------
Add NO_LIMIT=true to your bot's .pd file

----------------------------------------------------------------------
 Standard Error and Output Log Files
----------------------------------------------------------------------
Logs are stored in the working logs directory for Poker Academy Pro.
The streams are sent to error.log and output.log files.

The working logs directory's location is dependant on the OS. You
can determine the working logs directory by examining GameInfo.LOG_DIR
at runtime.

On Windows 7,8:

   C:\Users\YOU_USER_NAME\AppData\Roaming\PokerGenius\logs
   
On Mac OS X:

   ~/Library/Preferences/pokergenius/logs
 