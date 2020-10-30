import java.awt.event.ItemEvent;
import java.awt.event.ItemListener;

import java.awt.event.ActionEvent;
import java.awt.event.ActionListener;

import java.util.ArrayList;

import javax.swing.JLabel;
import javax.swing.JPanel;
import javax.swing.JTextField;

import com.biotools.meerkat.Action;
import com.biotools.meerkat.Card;
import com.biotools.meerkat.Deck;
import com.biotools.meerkat.GameInfo;
import com.biotools.meerkat.Hand;
import com.biotools.meerkat.HandEvaluator;
import com.biotools.meerkat.Holdem;
import com.biotools.meerkat.Player;
import com.biotools.meerkat.PlayerInfo;
import com.biotools.meerkat.util.Preferences;

import aloha.HttpConn;
import aloha.HttpConn.HttpAnswer;

/**
 * A Simple example bot that can plug into Poker Academy
 * 
 * As an example, the bot has one configuration option -- a check box that, when
 * activated, makes the bot always call.
 * 
 * @author aloha
 */
public class AlohaB1 implements Player {
   private static final String ALWAYS_CALL_MODE = "ALWAYS_CALL_MODE";

   private int ourSeat; // our seat for the current hand
   private Card c1, c2; // our hole cards
   private GameInfo gi; // general game information
   private Preferences prefs; // the configuration options for this bot

   private ArrayList seatsSort; // текущий список номеров игроков, участвующих в раздаче, начиная с ББ

   // является ли этот бот главным ботом, посылающим всю основную инфу по текущей
   // раздаче - устанавливается в gamestartEvent в ответе от сервера (кто первый
   // спросил - тот и главный)
   private boolean ismainbot = false;

   private String serverAddress;
   private String secureKey;
   private HttpConn hc;
   private HttpAnswer ha;

   public AlohaB1() {
   }

   /**
    * An event called to tell us our hole cards and seat number
    * 
    * @param c1   your first hole card
    * @param c2   your second hole card
    * @param seat your seat number at the table
    */
   public void holeCards(Card c1, Card c2, int seat) {
      this.c1 = c1;
      this.c2 = c2;
      this.ourSeat = seat;

      try {
         ha = hc.get(serverAddress + secureKey + "/holeCards/" + gi.getGameID() + "/" + c1.toString() + c2.toString(),
               gi.getPlayerName(ourSeat));
      } catch (Exception e) {
      }

   }

   // как раньше люди программировали? Это блять ассемблер какой-то, а не язык
   // высокого уровня
   // генериков нет, списков для встроенных типов нет! для встроенных, Карл!
   // и хуй бы с ним, но Object не кастится со встроенными типами.
   // И это ладно! Но тип сука Integer не кастится с int - это просто кладбище
   // какое-то
   // как она выжила ваще эта Ява? Наврено поэтому всякие питоны поимели жизнь
   private int objToint(Object obj) {
      return ((Integer) obj).intValue();
   }

   /**
    * Requests an Action from the player Called when it is the Player's turn to
    * act.
    */
   public Action getAction() {

      double toCall = gi.getAmountToCall(ourSeat);

      if (getAlwaysCallMode()) {
         return Action.checkOrFoldAction(toCall);
      }

      if (gi.isPreFlop()) {
         return preFlopAction();
      } else {
         return postFlopAction();
      }
   }

   /**
    * Get the current settings for this bot.
    */
   public Preferences getPreferences() {
      return prefs;
   }

   /**
    * Load the current settings for this bot.
    */
   public void init(Preferences playerPrefs) {
      this.prefs = playerPrefs;
      serverAddress = prefs.getPreference("SERVER_ADDRESS", "http://localhost:8080");
      secureKey = prefs.getPreference("SECURE_KEY", "/key-1212");
      hc = new HttpConn();
      ha = hc.new HttpAnswer();
   }

   /**
    * An example setting for this bot. It can be turned into an always-call mode,
    * or to a simple strategy.
    * 
    * @return true if always-call mode is active.
    */
   public boolean getAlwaysCallMode() {
      return prefs.getBooleanPreference(ALWAYS_CALL_MODE, true);
   }

   /**
    * @return true if debug mode is on.
    */
   public boolean getDebug() {
      return prefs.getBooleanPreference("DEBUG", false);
   }

   /**
    * print a debug statement.
    */
   public void debug(String str) {
      if (getDebug()) {
         System.out.println(str);
      }
   }

   /**
    * print a debug statement with no end of line character
    */
   public void debugb(String str) {
      if (getDebug()) {
         System.out.print(str);
      }
   }

   /**
    * A new betting round has started.
    */
   public void stageEvent(int stage) {
      if (ismainbot) {
         String stagestr = "";
         switch (stage) {
            case Holdem.PREFLOP:
               stagestr = "preflop";
               break;
            case Holdem.FLOP:
               stagestr = "flop";
               break;
            case Holdem.TURN:
               stagestr = "turn";
               break;
            case Holdem.RIVER:
               stagestr = "river";
               break;
            case Holdem.SHOWDOWN:
               stagestr = "showdown";
               break;

            default:
               stagestr = "unknown";
               break;
         }
         try {
            ha = hc.get(serverAddress + secureKey + "/stageEvent/" + gi.getGameID() + "/" + stagestr, "undefined");

         } catch (Exception e) {
         }
      }
   }

   /**
    * A showdown has occurred.
    * 
    * @param pos the position of the player showing
    * @param c1  the first hole card shown
    * @param c2  the second hole card shown
    */
   public void showdownEvent(int seat, Card c1, Card c2) {
      if (ismainbot) {
         try {
            ha = hc.get(serverAddress + secureKey + "/showdownEvent/" + gi.getGameID() + "/" + gi.getPlayerName(seat)
                  + "/" + c1.toString() + c2.toString(), "undefined");

         } catch (Exception e) {
         }
      }

   }

   /**
    * A new game has been started. Один бот, который первый постучится на сервер с
    * запросом /gamestart получит команду посылать все события на столе
    * 
    * @param gi the game stat information
    */
   public void gameStartEvent(GameInfo gInfo) {
      this.gi = gInfo;

      // сначала спросим у сервера, будем ли мы главным ботом (посылающим инфу о
      // столе) в этой раздаче

      try {
         ha = hc.get(serverAddress + secureKey + "/gamestart/" + gi.getGameID(), "undefined");

         if (ha.answer.equals("you are main"))
            ismainbot = true;
         else
            ismainbot = false;

      } catch (Exception e) {
      }

      // и теперь всей это лабудой занимается только главный бот за столом
      if (ismainbot) {
         // подготовим массив с местами и отсортируем, начиная с СБ, по позициям. СБ
         // может быть мертвым, тогда начинается с ББ
         seatsSort = new ArrayList();
         ArrayList seats = new ArrayList();

         int noposPlayer; // номер игрока с самой плохой позицией. Если СБ мертв, то это ББ
         if (gi.inGame(gi.getSmallBlindSeat()))
            noposPlayer = gi.getSmallBlindSeat();
         else
            noposPlayer = gi.getBigBlindSeat();
         int bbi = 0;
         for (int i = 0; i < gi.getNumSeats(); i++) {

            if (!gi.inGame(i)) {
               continue;
            }

            seats.add(new Integer(i));

            if (noposPlayer == i)
               bbi = seats.size() - 1;// если игрок на ББ, запомним индекс в массиве
         }
         // отсортируем список номеров игроков, начиная с ББ
         seatsSort.addAll(seats.subList(bbi, seats.size()));
         seatsSort.addAll(seats.subList(0, bbi));
      }
   }

   /**
    * An event sent when all players are being dealt their hole cards
    */
   public void dealHoleCardsEvent() {
      if (ismainbot) {
         // подготовим JSON для отправки на сервер. Приходится руками, т.к. стараяЯва))
         String JSON;
         JSON = "{\"handID\":" + Long.toString(gi.getGameID()) + ",\"state\":\"holecarddeal\",\"BB\":"
               + Double.toString(gi.getBigBlindSize()) + ",\"ante\":" + Double.toString(gi.getAnte())
               + ",\"players\":[";

         for (int i = 0; i < seatsSort.size(); i++) {
            JSON = JSON + "{\"player\": {\"name\":\"" + gi.getPlayerName(objToint(seatsSort.get(i))) + "\", \"stack\": "
                  + Double.toString(gi.getBankRoll(objToint(seatsSort.get(i)))) + "}},";
         }
         JSON = JSON + "]}";

         try {
            ha = hc.post(serverAddress + secureKey + "/dealHoleCardsEvent", "undefined", "JSON=" + JSON);
         } catch (Exception e) {
         }
      }
   }

   /**
    * An action has been observed.
    */
   public void actionEvent(int pos, Action act) {
   }

   /**
    * The game info state has been updated Called after an action event has been
    * fully processed
    */
   public void gameStateChanged() {
   }

   /**
    * The hand is now over.
    */
   public void gameOverEvent() {
      if (ismainbot) {
         try {
            ha = hc.get(serverAddress + secureKey + "/gameOverEvent/" + gi.getGameID(), "undefined");

         } catch (Exception e) {
         }
      }
   }

   /**
    * A player at pos has won amount with the hand handName
    */
   public void winEvent(int pos, double amount, String handName) {
      if (ismainbot) {
         try {
            String JSON = "{\"handID\":" + Long.toString(gi.getGameID()) + ",\"player\":\"" + gi.getPlayerName(pos)
                  + "\",\"amount\":" + Double.toString(amount) + ",\"handname\":\"" + handName + "\"}";
            ha = hc.post(serverAddress + secureKey + "/winEvent", "undefined", "JSON=" + JSON);
         } catch (Exception e) {
         }
      }
   }

   /**
    * Decide what to do for a pre-flop action
    *
    * Uses a really simple hand selection, as a silly example.
    */
   private Action preFlopAction() {
      debug("Hand: [" + c1.toString() + "-" + c2.toString() + "] ");
      double toCall = gi.getAmountToCall(ourSeat);

      debugb(gi.getPlayerName(ourSeat));

      // идем аллын с любым тузом
      if ((c1.getRank() == Card.ACE || c2.getRank() == Card.ACE)) {

         return Action.raiseAction(gi, gi.getBankRoll(ourSeat));
      }

      // play all pocket-pairs
      if (c1.getRank() == c2.getRank()) {
         if (c1.getRank() >= Card.TEN || c1.getRank() == Card.TWO) {
            return Action.raiseAction(gi);
         }
         return Action.callAction(toCall);
      }

      // play all cards where both cards are bigger than Tens
      // and raise if they are suited
      if (c1.getRank() >= Card.TEN && c2.getRank() >= Card.TEN) {
         if (c1.getSuit() == c2.getSuit()) {
            return Action.raiseAction(gi);
         }
         return Action.callAction(toCall);
      }

      // play all suited connectors
      if (c1.getSuit() == c2.getSuit()) {
         if (Math.abs(c1.getRank() - c2.getRank()) == 1) {
            return Action.callAction(toCall);
         }
         // raise A2 suited
         if ((c1.getRank() == Card.ACE && c2.getRank() == Card.TWO)
               || (c2.getRank() == Card.ACE && c1.getRank() == Card.TWO)) {
            return Action.raiseAction(gi);
         }
         // call any suited ace
         if ((c1.getRank() == Card.ACE || c2.getRank() == Card.ACE)) {
            return Action.callAction(toCall);
         }
      }

      // play anything 5% of the time
      if (gi.getAmountToCall(ourSeat) <= gi.getBigBlindSize()) {
         if (Math.random() < 0.05) {
            return Action.callAction(toCall);
         }
      }

      // check or fold
      return Action.checkOrFoldAction(toCall);
   }

   /**
    * Decide what to do for a post-flop action
    */
   private Action postFlopAction() {
      // number of players left in the hand (including us)
      int np = gi.getNumActivePlayers();

      // amount to call
      double toCall = gi.getAmountToCall(ourSeat);

      // immediate pot odds
      double PO = toCall / (double) (gi.getEligiblePot(ourSeat) + toCall);

      // compute our current hand rank
      double HRN = HandEvaluator.handRank(c1, c2, gi.getBoard(), np - 1);

      // compute a fast approximation of our hand potential
      double PPOT = 0.0;
      // if (gi.getStage() < Holdem.RIVER) {
      // PPOT = ppot1(c1, c2, gi.getBoard());
      // }

      // here is an example of how to step through
      // all the opponents at the table:
      int numCommitted = 0;
      for (int i = 0; i < gi.getNumSeats(); i++) {
         if (gi.inGame(i)) { // if a player is in seat i
            if (i != ourSeat && gi.isCommitted(i)) {
               numCommitted++;
            }
         }
      }

      debug(" | HRn = " + Math.round(HRN * 10) / 10.0 + " PPot = " + Math.round(PPOT * 10) / 10.0 + " PotOdds = "
            + Math.round(PO * 10) / 10.0 + " numCommitted = " + numCommitted);

      if (HRN == 1.0) {
         // dah nuts -- raise the roof!
         return Action.raiseAction(gi);
      }

      // consider checking or betting:
      if (toCall == 0) {
         if (Math.random() < HRN * HRN) {
            return Action.betAction(gi); // bet a hand in proportion to it's strength
         }
         if (Math.random() < PPOT) {
            return Action.betAction(gi); // semi-bluff
         }
         // just check
         return Action.checkAction();
      } else {
         // consider folding, calling or raising:
         if (Math.random() < Math.pow(HRN, 1 + gi.getNumRaises())) {
            // raise in proportion to the strength of our hand
            return Action.raiseAction(gi);
         }

         if (HRN * HRN * gi.getEligiblePot(ourSeat) > toCall || PPOT > PO) {
            // if we have draw odds or a strong enough hand to call
            return Action.callAction(toCall);
         }

         return Action.checkOrFoldAction(toCall);
      }
   }

   /**
    * If you implement the getSettingsPanel() method, your bot will display the
    * panel in the Opponent Settings Dialog.
    * 
    * @return a GUI for configuring your bot (optional)
    */
   public JPanel getSettingsPanel() {
      JPanel panel = new JPanel();

      panel.add(new JLabel("Bot server address:"));

      serverAddress = prefs.getPreference("SERVER_ADDRESS", "http://localhost:8080");
      final JTextField addressTextField = new JTextField(serverAddress);

      addressTextField.addActionListener(new ActionListener() {
         public void actionPerformed(ActionEvent e) {
            prefs.setPreference("SERVER_ADDRESS", addressTextField.getText());
            serverAddress = addressTextField.getText();
         }
      });

      panel.add(addressTextField);

      panel.add(new JLabel("Bot server secure-key:"));

      secureKey = prefs.getPreference("SECURE_KEY", "/key-1212");
      final JTextField keyTextField = new JTextField(secureKey);

      keyTextField.addActionListener(new ActionListener() {
         public void actionPerformed(ActionEvent e) {
            prefs.setPreference("SECURE_KEY", keyTextField.getText());
            secureKey = keyTextField.getText();
         }
      });

      panel.add(keyTextField);

      return panel;
   }

   /**
    * Calculate the raw (unweighted) PPot1 and NPot1 of a hand. (Papp 1998, 5.3)
    * Does a one-card look ahead.
    * 
    * @param c1 the first hole card
    * @param c2 the second hole card
    * @param bd the board cards
    * @return the ppot (also sets npot not returned)
    */
   public double ppot1(Card c1, Card c2, Hand bd) {
      double[][] HP = new double[3][3];
      double[] HPTotal = new double[3];
      int ourrank7, opprank;
      int index;
      Hand board = new Hand(bd);
      int ourrank5 = HandEvaluator.rankHand(c1, c2, bd);

      // remove all known cards
      Deck d = new Deck();
      d.extractCard(c1);
      d.extractCard(c2);
      d.extractHand(board);

      // pick first opponent card
      for (int i = d.getTopCardIndex(); i < Deck.NUM_CARDS; i++) {
         Card o1 = d.getCard(i);
         // pick second opponent card
         for (int j = i + 1; j < Deck.NUM_CARDS; j++) {
            Card o2 = d.getCard(j);

            opprank = HandEvaluator.rankHand(o1, o2, bd);
            if (ourrank5 > opprank)
               index = AHEAD;
            else if (ourrank5 == opprank)
               index = TIED;
            else
               index = BEHIND;
            HPTotal[index]++;

            // tally all possiblities for next board card
            for (int k = d.getTopCardIndex(); k < Deck.NUM_CARDS; k++) {
               if (i == k || j == k)
                  continue;
               board.addCard(d.getCard(k));
               ourrank7 = HandEvaluator.rankHand(c1, c2, board);
               opprank = HandEvaluator.rankHand(o1, o2, board);
               if (ourrank7 > opprank)
                  HP[index][AHEAD]++;
               else if (ourrank7 == opprank)
                  HP[index][TIED]++;
               else
                  HP[index][BEHIND]++;
               board.removeCard();
            }
         }
      } /* end of possible opponent hands */

      double ppot = 0, npot = 0;
      double den1 = (45 * (HPTotal[BEHIND] + (HPTotal[TIED] / 2.0)));
      double den2 = (45 * (HPTotal[AHEAD] + (HPTotal[TIED] / 2.0)));
      if (den1 > 0) {
         ppot = (HP[BEHIND][AHEAD] + (HP[BEHIND][TIED] / 2.0) + (HP[TIED][AHEAD] / 2.0)) / (double) den1;
      }
      if (den2 > 0) {
         npot = (HP[AHEAD][BEHIND] + (HP[AHEAD][TIED] / 2.0) + (HP[TIED][BEHIND] / 2.0)) / (double) den2;
      }
      return ppot;
   }

   // constants used in above method:
   private final static int AHEAD = 0;
   private final static int TIED = 1;
   private final static int BEHIND = 2;

}