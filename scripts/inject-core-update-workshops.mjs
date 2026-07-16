#!/usr/bin/env node
/**
 * Turn each LEVEL 01–12 Update method into a paste-in-order workshop.
 *
 * Code is extracted from the real Go file.  The prose table deliberately maps
 * one explanation to every blank-line-delimited logic block; a source edit
 * therefore fails generation until the lesson is updated too.
 */
import { readFileSync, writeFileSync } from "node:fs";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("..", import.meta.url));
const markerStart = "<!-- core-update-workshop:start -->";
const markerEnd = "<!-- core-update-workshop:end -->";

const lessons = [
  {
    slug: "tap-target", source: "games/core/tap-target/main.go", level: "01",
    ja: { title: "入力から得点まで、Updateを5つにほどく", lead: "最初のUpdateは、ゲーム開始・時間切れ・タイマー・当たり判定を順番に判断します。短い処理でも、上から読むだけで一つの物語になります。", blocks: [
      ["入口で入力を一度だけ読む", "マウスとタッチを px・py・pressed の3値へそろえます。複数の戻り値を一行で受け取れるのはGoの読みやすいところです。以後の判断は、このtickで押されたかだけを見ます。", "左辺の3つと右辺の3つが同じ順番で対応します。"],
      ["始まる前なら、開始だけを扱う", "! は「ではない」です。まだ始まっていない間は、押された時だけ started を true にして早くreturnします。ここで戻るため、同じ入力が得点にも使われません。", "早期returnは、今関係ない処理を読まなくてよくする道案内です。"],
      ["時間切れなら、リセットだけを扱う", "残りtickが0以下なら終了画面の状態です。押されたら得点・時間・大きさを初期状態へ戻します。複数の代入が『新しい一回』を作ります。", "<= は0ちょうどだけでなく、念のため負の値も含めます。"],
      ["遊んでいる時だけ時計を進める", "ここまでreturnされなかった時だけ、残り時間を1 tick減らします。g.framesLeft-- は g.framesLeft = g.framesLeft - 1 の短い書き方です。", "コードの置き場所がルールになります。開始前にはこの行へ届きません。"],
      ["押した瞬間に、当たりを一度だけ採点する", "pressed の時だけ円の判定を呼び、当たれば score を1増やします。math.Maxで小さくなりすぎるのを防ぎ、次の的を動かします。最後の nil は『問題なくこのtickを終えた』です。", "ifを入れ子にすると『押した、しかも当たった』をそのまま読めます。"],
    ]},
    en: { title: "Unpack Update from input to score", lead: "This first Update makes four decisions in order: start, timeout, timer, and hit. Even a small function can read like a story from top to bottom.", blocks: [
      ["Read input once at the entrance", "Mouse and touch become px, py, and pressed. Go can return several values cleanly, and every later rule uses the same snapshot for this tick.", "The three names on the left match the three returned values in order."],
      ["Before play, handle only Start", "! means “not.” While the game has not started, a press flips started to true and an early return prevents that same press from scoring.", "An early return is a signpost: skip rules that do not belong to this state."],
      ["After timeout, handle only Reset", "At zero ticks the run is over. A press restores score, time, and target state so the next run is genuinely fresh.", "<= safely covers zero and any accidental negative value."],
      ["Advance the clock only during play", "Only a tick that passed both early returns reaches framesLeft--. It is shorthand for subtracting one and storing the result.", "Placement is part of the rule: the timer cannot move on the start screen."],
      ["Score one decision on the press", "Only a fresh press runs the circle test. A hit increments score, clamps the shrinking radius with math.Max, and moves the target. nil means this tick ended normally.", "Nested if statements read as “pressed, and then hit.”"],
    ]},
  },
  {
    slug: "timing-meter", source: "games/core/timing-meter/main.go", level: "02",
    ja: { title: "同じボタンの意味を、状態で切り替える", lead: "押すたびに止める／再開するが切り替わります。bool、剰余、早期returnを使うと、状態機械を日本語の手順のように書けます。", blocks: [
      ["押されたtickだけ、モードを変える", "justPressedは押しっぱなしを毎tick数えません。stoppedがtrueなら次ラウンド、falseなら停止と採点です。elseによって二つが同時に起きないこともコードが保証します。", "%2 は偶数・奇数を見分ける定番です。速度の符号を反転すると向きが変わります。"],
      ["止まっている間は、位置を変えない", "採点後はここでreturnするので、下の移動式に届きません。Drawは止まったmarkerXを何度描いてもよく、停止という事実はUpdateが守ります。", "状態を変えないtickも、立派なUpdateの結果です。"],
      ["動作中だけ、移動と反射を計算する", "純粋関数が新しい位置と速度を2つ返します。左右の端を越えた時は速度の符号も変わります。最後にnilで正常終了します。", "入力と出力がはっきりした関数は、単体テストもしやすくなります。"],
    ]},
    en: { title: "Let state change what one button means", lead: "Each press alternates between stop and restart. A bool, modulo, and an early return turn that into a state machine that reads like instructions.", blocks: [
      ["Change mode only on a fresh press", "justPressed avoids counting a held button every tick. stopped chooses either next-round setup or scoring; else guarantees both cannot happen together.", "%2 is the classic even/odd check. Negating speed reverses direction."],
      ["While stopped, do not move", "The early return keeps the movement equation unreachable after scoring. Draw may present the same marker many times; Update owns the fact that it is stopped.", "A tick with no position change is still a valid Update result."],
      ["Move and bounce only while running", "A pure helper returns both the next position and speed. Crossing an edge changes the speed sign. nil ends the tick normally.", "Clear inputs and outputs make a helper pleasant to unit-test."],
    ]},
  },
  {
    slug: "catch-stars", source: "games/core/catch-stars/main.go", level: "03",
    ja: { title: "増える星を、順番に生み・動かし・選び直す", lead: "LEVEL 03ではsliceが登場します。星が何個でも、同じforを一度書けば全員へ同じルールを適用できます。", blocks: [
      ["終了中はリスタート以外を止める", "gameOverなら新しいgameの値を丸ごと代入し、すぐ戻ります。古い星や乱数状態まで初期化されるので、再開漏れを減らせます。", "*g = *newGame() は、今の箱の中身を新しい初期値で置き換えます。"],
      ["先にカゴを動かす", "当たり判定より前に入力を位置へ反映します。そのtickで動いた先に星があればキャッチできる、という自然な操作感になります。", "Update内の順番は、プレイヤーが感じる因果関係です。"],
      ["一定間隔で新しい星をappendする", "frameを増やし、45で割り切れるtickだけstarを作ります。構造体リテラルは、x・y・speedを名前付きでひとまとめにできます。", "appendはsliceの末尾へ新しい要素を安全に足すGoの組み込み関数です。"],
      ["残す星の入れ物を用意し、一つずつ動かす", "next := g.stars[:0] は同じ領域を再利用する空sliceです。rangeで星を値として取り出し、落下後のsを判定します。", "『消す物を探す』より『残す物だけ集める』と考えると処理が読みやすくなります。"],
      ["結果ごとにcontinueし、最後に残った星を採用する", "switchでキャッチ・落下・継続を分けます。continueした星はnextへ追加されないため消え、どちらでもない星だけ残ります。", "switchは一つの結果に対する複数の道を、縦に読みやすく並べます。"],
    ]},
    en: { title: "Spawn, move, and filter a growing star list", lead: "LEVEL 03 introduces slices. One loop can apply the same rule to any number of stars.", blocks: [
      ["During game over, allow only restart", "Replacing *g with a fresh game resets old stars and random state together, then returns before gameplay runs.", "*g = *newGame() replaces the contents of the current game value."],
      ["Move the basket first", "Input changes position before collision checks, so a basket moved onto a star catches it in the same tick.", "Order inside Update becomes cause and effect the player can feel."],
      ["Append a star at a steady interval", "Increment frame and create a star only when it divides evenly by 45. A struct literal groups named x, y, and speed values.", "append is Go's built-in way to grow a slice."],
      ["Prepare a keep-list and move every star", "g.stars[:0] makes an empty slice that reuses storage. range visits each star, then its copied value is advanced.", "Thinking “keep survivors” is often clearer than deleting while moving forward."],
      ["Continue on outcomes; keep everything else", "switch separates caught, missed, and still falling. continue skips append, so those stars disappear; the rest enter next.", "switch lays several paths from one result out vertically."],
    ]},
  },
  {
    slug: "flappy", source: "game/main.go", level: "04",
    ja: { title: "一回の羽ばたきから、重力のある世界を作る", lead: "速度を先に変え、その速度で位置を変える。この順番を守るだけで、数字が滑らかなジャンプの弧になります。", blocks: [
      ["tick番号と入力を記録する", "frameを一つ増やし、このtickで押されたかを一度だけ読みます。以後は同じpressedを共有するため、処理の途中で入力が変わる心配がありません。", ":= は型をGoに推測してもらいながら新しい変数を作る記号です。"],
      ["終了中はリセットだけ", "ゲームオーバー後に押されたらresetを呼びます。returnで重力やパイプ処理を止めるため、結果画面の裏で世界が進みません。", "状態ごとに入口を分けると、バグの起きる組み合わせが減ります。"],
      ["開始前の待機アニメと最初の羽ばたき", "Sinで上下に揺らし、押されたらstartedと上向き速度を設定します。見た目用のbirdYもUpdateで決まり、Drawはその数字を読むだけです。", "小さな数式が動きに性格を与えるのが、ゲームプログラミングの楽しいところです。"],
      ["押した瞬間に速度を上向きへ", "位置を直接持ち上げず、velocityだけをflapSpeedにします。だから次の重力計算と自然につながります。", "操作は『結果の位置』ではなく『運動への力』を変えています。"],
      ["重力を足し、位置を進める", "純粋関数は現在位置・速度・加速度から次の位置と速度を返します。この2行分の物理を独立させると、画面なしでテストできます。", "多値代入なら、関連する二つの状態を同じ式の結果で更新できます。"],
      ["全パイプを左へ動かし、通過を一度だけ数える", "rangeのindexでslice内の本物の要素を書き換えます。scoredフラグは同じパイプから二重点を取らないための記憶です。", "bool一つが『もう済んだか』を覚える小さな状態機械になります。"],
      ["画面外のパイプを再利用する", "先頭を捨て、末尾の右へ新しいpipeをappendします。配列全体を作り直さず、流れる列を保てます。", "lenで最後の位置を求めるため、パイプ数を変えても同じコードが使えます。"],
      ["衝突をまとめ、最高点を保存する", "地面・天井・パイプのどれかならgameOverです。|| は『どれか一つ』、最後のifは今回の得点がbestを超えた時だけ記録します。", "複雑な条件も名前を付けると、英語の文章のように読めます。"],
    ]},
    en: { title: "Build a world with gravity from one flap", lead: "Change velocity first, then position. That order turns a few numbers into a smooth jump arc.", blocks: [
      ["Record the tick and one input snapshot", "Increment frame and read a fresh press once. Every later rule shares the same pressed value.", ":= creates a variable while Go infers its type."],
      ["During game over, reset only", "A press calls reset, while return prevents gravity and pipes from continuing behind the result screen.", "State-specific entrances reduce impossible combinations."],
      ["Idle motion and the first flap", "Sin adds a gentle bob. A press sets started and upward velocity. Even presentation motion is decided in Update; Draw only reads birdY.", "A tiny formula giving motion personality is one joy of game programming."],
      ["On a press, change upward velocity", "Do not teleport position. Change velocity, which naturally feeds the gravity rule next.", "Input changes the motion, not its final answer."],
      ["Add gravity, then advance position", "A pure helper returns next position and velocity from current state and acceleration, so this physics can be tested without a window.", "Multiple assignment updates related state from one result."],
      ["Move every pipe and score each once", "An index from range mutates the real slice element. scored remembers that one pipe has already paid out.", "One bool can be a tiny state machine."],
      ["Recycle the off-screen pipe", "Drop the first pipe and append a new one after the last, preserving a flowing queue without rebuilding everything.", "len keeps this code useful if the pipe count changes."],
      ["Combine collisions and save the best", "Ground, ceiling, or pipe sets gameOver. Named booleans make the || condition readable; another if records only a new best.", "Naming conditions turns math into a sentence."],
    ]},
  },
  {
    slug: "pong", source: "games/core/pong/main.go", level: "05",
    ja: { title: "大きなUpdateを、仕事ごとのメソッドへ分ける", lead: "PongのUpdateは短く見えます。難しい処理を隠したのではなく、移動・AI・反射という意味のある名前へ分けたからです。", blocks: [
      ["一つのtickを、三つの動詞で読む", "プレイヤー、CPU、ボールの順に状態を更新します。詳細は各メソッドにあり、Updateからはゲーム全体の流れがひと目で見えます。", "Goの小さなメソッドは、コメントより正確な『動く見出し』になります。"],
      ["画面外へ出た向きで得点を分ける", "ExitScoreの結果をswitchで受け、上ならプレイヤー、下ならCPUへ加点します。serveの引数1と-1が次のボール方向です。", "同じ型の処理をcaseで並べると、追加や変更が簡単です。"],
    ]},
    en: { title: "Split a large Update into named jobs", lead: "Pong's Update looks short because movement, AI, and reflection have meaningful method names—not because the hard parts disappeared.", blocks: [
      ["Read one tick as three verbs", "Update player, CPU, then ball. Details live in small methods while Update shows the whole game's order at a glance.", "Small Go methods are executable headings, more precise than comments."],
      ["Score by the side the ball exited", "switch maps ExitScore to player or CPU points. serve receives 1 or -1 for the next direction.", "Cases keep parallel outcomes easy to extend."],
    ]},
  },
  {
    slug: "breakout", source: "games/core/breakout/main.go", level: "06",
    ja: { title: "ボール、命、ブロックを小さな規則でつなぐ", lead: "複数のブロックがあってもUpdateの骨格は短くできます。まず世界を進め、その結果として失敗やクリアを判定します。", blocks: [
      ["操作・反射・破壊を順に進める", "パドルを動かしてから壁、最後にブロックとの衝突を処理します。一つ前の結果を次の仕事が読むため、順序がはっきりします。", "メソッド分割は、どこを直せばよいかも教えてくれます。"],
      ["落下したら命を一つ使う", "SpendLifeは新しいlivesとgameOverを返します。最後の命なら全初期化、まだあればserveだけに分岐します。", "計算を純粋関数へ出すと、境界値1と0をテストしやすくなります。"],
      ["残り0個をクリア条件にする", "aliveBrickCountが0なら新しいゲームへ置き換えます。『全部壊した』を数として表すため、ブロックの総数に依存しません。", "状態を数える関数名が、そのまま勝利条件の説明になります。"],
    ]},
    en: { title: "Connect ball, lives, and bricks with small rules", lead: "Even with many bricks, Update stays readable: advance the world first, then react to failure or victory.", blocks: [
      ["Advance control, bounce, then destruction", "Move the paddle, resolve world bounces, then brick hits. Each job reads the previous result in an explicit order.", "Method boundaries also tell you where a future change belongs."],
      ["Spend one life when the ball falls", "SpendLife returns both new lives and gameOver. The final life resets everything; otherwise only serve again.", "A pure helper makes boundary values such as one and zero easy to test."],
      ["Use zero remaining bricks as victory", "When aliveBrickCount is zero, replace the game. The rule works regardless of the original brick count.", "A counting method can read exactly like the win condition."],
    ]},
  },
  {
    slug: "snake", source: "games/core/snake/main.go", level: "07",
    ja: { title: "毎tick読む入力と、時々だけ進むマス移動を分ける", lead: "方向キーはいつでも受け付けたい一方、ヘビは一定間隔だけ1マス進みます。この二つの時間をUpdate内で分けます。", blocks: [
      ["終了中は再開だけを受け付ける", "gameOverなら新しいgameへ置き換え、通常入力や移動を止めます。returnによって死んだヘビが裏で動きません。", "ガード節は『この状態ではここまで』を最初に宣言します。"],
      ["方向入力は毎tick読む", "移動tickでなくてもreadInputを呼ぶため、素早いキー操作を取りこぼしにくくなります。実際の1マス移動はまだ行いません。", "入力の記憶と、世界の進行を別々に考えます。"],
      ["必要なtickだけstepSnakeする", "scoreから待ち時間を求め、frameが割り切れない間はreturnします。割り切れた時だけ頭・体・餌・衝突を一段進めます。", "% は周期処理を短く表し、速さ調整も一つの関数に閉じ込めます。"],
    ]},
    en: { title: "Separate always-listening input from occasional grid steps", lead: "Direction keys should be remembered every tick, while the snake advances one cell only at an interval.", blocks: [
      ["During game over, accept restart only", "Replace the game and return so a dead snake cannot keep moving behind the overlay.", "A guard clause declares “this state stops here.”"],
      ["Read direction every tick", "Calling readInput even between movement ticks avoids missing a quick key press. No grid movement happens yet.", "Input memory and world progression are separate ideas."],
      ["Call stepSnake only on scheduled ticks", "Derive wait from score and return until frame divides evenly. Only then advance head, body, food, and collision one step.", "% expresses periodic work compactly, while a helper owns speed tuning."],
    ]},
  },
  {
    slug: "space-shooter", source: "games/core/space-shooter/main.go", level: "08",
    ja: { title: "三種類のリストが動く、シューティングの管制塔", lead: "自機、敵、自弾、敵弾を一つのUpdateが順番に指揮します。長くなっても、各ブロックの入力と結果を追えば迷いません。", blocks: [
      ["終了状態を先に閉じる", "ゲームオーバーならリスタートだけを見て戻ります。通常世界へ進む入口が一つに保たれます。", "長い関数ほど、最初のガード節が安全な範囲をはっきりさせます。"],
      ["共通時計と無敵時間を進める", "frameを増やし、invincibleが正の間だけ1減らします。> 0 のガードがあるので負の値になりません。", "カウンターは『あと何tick続くか』を表す最小の状態です。"],
      ["背景の星を循環させる", "rangeのindexで全ての星を動かし、下へ抜けた星だけ上へ戻します。装飾もUpdateで位置が決まり、Drawはその位置を映します。", "同じ配列を循環利用すると、毎回作って捨てる必要がありません。"],
      ["自機の入力と発射を処理する", "移動と射撃を名前付きメソッドへ分けます。Updateのこの2行だけで、プレイヤー操作の順序が読めます。", "『一つのメソッドに一つの動詞』は、コードを組み替えやすくします。"],
      ["二種類の弾を同じ形で動かす", "自弾と敵弾は別sliceですが、x += vx、y += vyという同じ規則を持ちます。役割は分け、運動の考え方はそろえます。", "+= は今の値へ速度を足して再代入する短い書き方です。"],
      ["敵ごとに移動し、狙って撃つ", "&g.enemies[i] でslice内の敵そのものを指し、再装填を減らします。射撃時は自機との差から速度を求め、新しい弾をappendします。", "ポインタを使う理由は、コピーでなく元の敵を書き換えるためです。"],
      ["自弾と敵を後ろから調べる", "削除しながらsliceを走査するためindexを末尾から減らします。命中した敵と弾を切り取り、得点を加えます。", "後ろから消せば、まだ調べていないindexがずれません。"],
      ["敵弾を片付け、当たればhurtする", "画面外なら先に削除してcontinue。残った弾だけ無敵時間と衝突を調べます。条件の安いものから落とすと読みやすく効率的です。", "continueは『この弾の話は終わり、次へ』を表します。"],
      ["突破した敵もダメージへ変える", "敵が画面下へ出たらsliceから消し、hurtを呼びます。弾の衝突と敵の突破が同じダメージ処理へ合流します。", "共通処理を一か所にすると、ライフ減少の仕様がずれません。"],
      ["敵が0なら次のwaveを作る", "lenは現在の敵数です。0になった瞬間にspawnWaveを呼び、次の遊びを用意してUpdateを終えます。", "『空になった』というデータが、そのまま場面転換の合図になります。"],
    ]},
    en: { title: "Use Update as air-traffic control for three lists", lead: "One Update directs player, enemies, player shots, and enemy shots in order. Follow each block's input and result and a long function stays understandable.", blocks: [
      ["Close the game-over state first", "Only restart is allowed before returning, leaving one clear entrance to the live world.", "A guard clause matters more as a function grows."],
      ["Advance the shared clock and invincibility", "Increment frame and decrement invincible only above zero, preventing a negative timer.", "A counter is the smallest state for “how many ticks remain.”"],
      ["Cycle background stars", "Mutate every star by index and wrap only those below the screen. Draw merely presents these Update-owned positions.", "Reusing objects avoids creating and discarding them every pass."],
      ["Handle player movement and shooting", "Named methods make the order of player actions visible in two lines.", "One verb per method makes code easy to rearrange."],
      ["Move two bullet lists with one idea", "Player and enemy bullets stay separate, but both apply x += vx and y += vy.", "+= adds to the current value and stores it back."],
      ["Move and fire each enemy", "A pointer to g.enemies[i] mutates the real enemy. Reload counts down; firing appends a velocity aimed from the position difference.", "The pointer matters because we want the original, not a copy."],
      ["Check player bullets backward", "Iterating from the end is safe while slicing out hit enemies and bullets, then adding score.", "Backward deletion does not shift indices you still need to visit."],
      ["Clean enemy bullets, then call hurt", "Remove off-screen shots and continue. Only survivors need invincibility and collision checks.", "continue says “this bullet is finished; visit the next.”"],
      ["Turn escaped enemies into damage", "Remove enemies below the screen and funnel the result through the same hurt method.", "One shared damage path prevents rule drift."],
      ["Spawn a wave when the list is empty", "len gives the current enemy count; zero becomes the signal to prepare the next challenge.", "Data itself can be a clean scene-transition signal."],
    ]},
  },
  {
    slug: "sokoban", source: "games/core/sokoban/main.go", level: "09",
    ja: { title: "一手を受け付ける前に、優先順位を決める", lead: "リセット、取り消し、クリア後、移動アニメ、次の入力。倉庫番は『今できること』を上から順に絞る教材です。", blocks: [
      ["最優先の操作をガード節で並べる", "Rはいつでも全初期化、Z/Backspaceは移動中でない時だけundo、クリア後は再開だけです。それぞれreturnするので、一つのtickに複数の命令が混ざりません。", "条件を優先順に上から並べると、仕様書とコードの順番が一致します。"],
      ["マス間を動いている途中は補間だけ", "moving中はprogressを進め、完了時だけ論理座標と箱位置を確定します。その間は新しい移動入力を読まず、半端な位置から次の手が始まりません。", "一時的なアニメ状態と、確定したゲーム状態を分けています。"],
      ["止まっている時だけ一手を読む", "readMoveの結果がゼロでなければmoveへ渡します。move側が壁や箱を調べ、可能な手だけmoving状態を始めます。", "ゼロ値が『入力なし』として自然に使えるのもGoの良さです。"],
    ]},
    en: { title: "Decide priority before accepting a move", lead: "Reset, undo, cleared, animation, next input: Sokoban narrows what is currently legal from top to bottom.", blocks: [
      ["Place highest-priority commands as guards", "R always resets; undo works only while still; cleared accepts restart only. Each returns so one tick cannot mix commands.", "Ordering conditions by priority keeps specification and code aligned."],
      ["While between tiles, advance only the tween", "moving advances progress and commits logical player/box positions only on completion. New movement input is ignored until then.", "Temporary animation state stays separate from committed game state."],
      ["Read one move only while settled", "A nonzero direction goes to move, which checks walls and boxes before starting a transition.", "Go's zero value naturally represents “no direction.”"],
    ]},
  },
  {
    slug: "platformer", source: "games/core/platformer/main.go", level: "10",
    ja: { title: "横と縦を分けると、足場の衝突が解ける", lead: "2Dの衝突を一度に解こうとすると、どちらへ押し戻すか分かりません。横を動かして直し、次に縦を動かして直します。", blocks: [
      ["終了状態を閉じる", "クリアまたはゲームオーバーなら、再開以外の物理を止めます。|| はどちらか一方でもtrueなら成立します。", "最初に世界を止めることで、結果表示中の落下や得点を防ぎます。"],
      ["三つの操作を一度に読む", "readControlsがleft、right、jumpを返します。呼び出し側は入力機器の細部を知らず、ゲームの意味だけを受け取ります。", "多値returnは、関連する小さな情報を構造体なしで渡すのに便利です。"],
      ["入力から横速度、ジャンプから縦速度を作る", "純粋関数へ現在速度と入力を渡し、次の速度を受け取ります。地面を離れた合図も返るためonGroundを更新できます。", "物理式をEbitengineから分けると、数値だけでテストできます。"],
      ["横だけ動かし、壁から押し戻す", "xへvxを足した後、全platformと重なりを調べます。進行方向を見て左右どちらへ戻すか決め、速度を0にします。", "一軸ずつ解くことで、押し戻す方向が一意になります。"],
      ["縦だけ動かし、床と天井を分ける", "yへvyを足し、落下中なら床の上へ、上昇中なら天井の下へ戻します。床に着いた時だけonGround=trueです。", "速度の正負が、ぶつかった面を教えてくれます。"],
      ["コインを後ろから取り除く", "末尾から調べ、重なったcoinをsliceから切り取ってcollectedを増やします。複数枚が同じtickに重なっても安全です。", "削除を伴うslice走査では後ろ向きが定番です。"],
      ["失敗とゴールを別々の条件にする", "画面下へ落ちればgameOver、goalXを越えればclearedです。勝敗の事実だけをUpdateに保存します。", "Drawはこのboolを読んで表示を変えるだけで、勝敗を決めません。"],
      ["カメラを目標へ少しずつ近づける", "プレイヤーを画面左寄りに置くtargetを作り、差の9%だけcameraXを進めます。最後にclampしてステージ外を見せません。", "『差の一部を足す』だけで滑らかな追従が生まれます。"],
    ]},
    en: { title: "Solve platform collisions one axis at a time", lead: "If 2D collision is solved all at once, the push-out direction is ambiguous. Move and fix X, then move and fix Y.", blocks: [
      ["Close terminal states", "Cleared or game over allows restart only, preventing physics behind the result screen.", "Stopping the world first prevents hidden falls and scoring."],
      ["Read three controls together", "readControls returns left, right, and jump, hiding device details from the game rule.", "Multiple return values are handy for related small facts."],
      ["Turn input into horizontal and vertical velocity", "Pure helpers accept current velocity and inputs, returning next velocity and a left-ground signal.", "Ebitengine-free math can be tested with numbers alone."],
      ["Move X, then push out of walls", "After adding vx, inspect platforms, choose a side from velocity sign, and zero velocity.", "One axis at a time gives an unambiguous correction."],
      ["Move Y, then distinguish floor and ceiling", "Positive velocity lands on top; negative velocity moves below a ceiling. Only landing sets onGround.", "Velocity sign tells you which face was hit."],
      ["Remove collected coins backward", "Visit from the end, splice overlaps out, and increment collected safely even for multiple hits.", "Backward iteration is the standard slice-deletion pattern."],
      ["Store failure and goal as separate facts", "Falling sets gameOver; crossing goalX sets cleared. Draw only presents those bools.", "Draw never decides victory."],
      ["Ease the camera toward a target", "Move 9% of the remaining difference and clamp to stage bounds.", "Adding a fraction of the gap creates smooth follow motion."],
    ]},
  },
  {
    slug: "dungeon", source: "games/core/dungeon/main.go", level: "11",
    ja: { title: "敵AIの『考える→動く→当たる』を順番に読む", lead: "長いUpdateでも、プレイヤー処理の後に各敵を一体ずつ進めるだけです。ポインタ、距離、状態遷移が一つの小さなAIになります。", blocks: [
      ["終了状態を閉じる", "クリアかゲームオーバーなら再開だけを受け付けます。AIも攻撃時間も進まないため、結果画面が安定します。", "先頭のガード節が、生きている世界と止まった世界を分けます。"],
      ["プレイヤーを横、次に縦へ動かす", "readMoveのdx、dyを一軸ずつmovePlayerへ渡します。壁衝突を方向ごとに解けます。", "LEVEL 10の一軸分解を、見下ろし移動にも再利用しています。"],
      ["攻撃の残り時間を減らす", "attackが正の時だけ--します。0は攻撃していない、1以上は攻撃中／無敵中という時間付き状態です。", "一つの整数がアニメとルールの両方の時刻表になります。"],
      ["複数の入力機器をattackPressedへまとめる", "Space、X、マウス、上半分のタッチを一つのboolへ統合し、クールダウン0の時だけ18を入れます。", "入力方法が増えても、その後のゲームルールは一つのboolだけを見ます。"],
      ["プレイヤー中心を一度計算する", "矩形左上ではなく中心で距離や剣先を計算するため、+14した座標へ名前を付けます。", "繰り返す式へ名前を付けると、意味と修正場所が一つになります。"],
      ["敵sliceの本物を一体ずつ指す", "&g.enemies[i]で元の敵を更新し、aliveでない敵はcontinueします。死んだ敵へAIを走らせません。", "ポインタとcontinueが、必要な対象だけを扱うループを作ります。"],
      ["距離から巡回／追跡を選ぶ", "Hypotで距離を求め、EnemyModeへ今のstateと二つのしきい値を渡します。追い始める距離と諦める距離を分けることでガタつきを防ぎます。", "状態遷移を純粋関数にすると、距離の境界を表形式でテストできます。"],
      ["巡回と追跡で速度の作り方を変える", "state 0なら一定時間ごとに方向転換、それ以外ならプレイヤー方向の単位速度を求めます。distが0なら割り算を守ります。", "同じe.vx/e.vyへ合流させるため、後の移動処理はAIの種類を知りません。"],
      ["敵も横、次に縦へ動かす", "決まった速度を二回のmoveEnemyへ渡します。AIの判断と衝突解決を別の仕事に保ちます。", "考える処理と動かす処理を分けると、新しいAIを足しやすくなります。"],
      ["剣先との距離で敵を倒す", "向き×30で剣先を作り、攻撃の前半と距離35以内が重なった時だけalive=falseにします。", "見えない当たり判定も、座標と距離という同じ道具で作れます。"],
      ["接触ダメージを一度だけ処理する", "近く、かつ攻撃中でない時にlifeを減らし、attack=55をのけぞり兼無敵時間にします。0以下ならgameOverです。", "&& は全条件を満たす時だけ危険を成立させます。"],
      ["全滅と出口位置を勝利条件にする", "aliveCountが0で、出口の座標範囲に入った時だけclearedです。敵を無視して出口へ走ることはできません。", "複数の事実を&&で組み合わせると、物語の条件になります。"],
    ]},
    en: { title: "Read enemy AI as think, move, collide", lead: "Even a long Update handles the player, then advances one enemy at a time. Pointers, distance, and state transitions make a small AI.", blocks: [
      ["Close terminal states", "Only restart is accepted after clear or game over, freezing AI and attack timers.", "The guard divides the live world from the stopped world."],
      ["Move player X, then Y", "Pass dx and dy separately so wall collision resolves per axis.", "This reuses LEVEL 10's axis decomposition from a top-down view."],
      ["Count down attack time", "Only positive attack decrements. Zero means idle; positive values are a timed attack/invulnerability state.", "One integer can schedule animation and rules."],
      ["Merge devices into attackPressed", "Keyboard, mouse, and touch become one bool; only cooldown zero starts 18 ticks.", "More input methods do not complicate the game rule."],
      ["Calculate player center once", "Name the +14 center used by distance and sword-tip math.", "Naming a repeated expression gives it meaning and one edit point."],
      ["Point to each real enemy", "&g.enemies[i] mutates the original; continue skips dead enemies.", "Pointers plus continue focus the loop on relevant objects."],
      ["Choose patrol or chase from distance", "EnemyMode receives current state and separate enter/exit thresholds, avoiding jitter.", "A pure transition helper invites table-driven boundary tests."],
      ["Build velocity differently by state", "Patrol turns periodically; chase aims at the player, guarding distance zero.", "Both paths join at vx/vy, so movement need not know the AI kind."],
      ["Move enemy X, then Y", "Two moveEnemy calls keep AI decisions separate from collision resolution.", "Separating thinking from movement makes new AI easier."],
      ["Defeat an enemy at the sword tip", "Facing times 30 gives the tip; active attack time and distance set alive=false.", "Invisible hitboxes use the same coordinates and distance tools."],
      ["Apply contact damage once", "Near and not attacking decrements life, starts recovery, and may set gameOver.", "&& makes danger require every condition."],
      ["Combine defeat-all and exit into victory", "No living enemies plus the exit region sets cleared.", "Combining facts with && creates a story condition."],
    ]},
  },
  {
    slug: "bullet-hell", source: "games/core/bullet-hell/main.go", level: "12",
    ja: { title: "大量の弾も、一発ずつ同じ規則で進める", lead: "弾幕は難しそうに見えますが、生成する周期と一発の更新を組み合わせたものです。sliceの後ろ向き走査で、移動・命中・削除を一周にまとめます。", blocks: [
      ["終了状態を閉じる", "クリアまたはゲームオーバー中は再開だけです。弾幕生成も移動も止まり、同じ結果をDrawが安全に表示できます。", "最初のreturnが、残り全ての処理の前提を作ります。"],
      ["共通時計と無敵時間を進める", "frameは発射スケジュール、invincibleは被弾間隔に使います。正の間だけ減らすので0で安定します。", "異なる仕組みでも、tickカウンターという同じ部品を再利用できます。"],
      ["プレイヤー位置を先に確定する", "弾の発射や敵弾の狙いより前にmovePlayerを呼びます。このtickの最新位置が後続の判定に使われます。", "順序を選ぶことが操作感を選ぶことになります。"],
      ["7tickごとに自弾を2発appendする", "%7==0を発射の時計にし、左右の砲口からbulletを追加します。構造体リテラル一つが位置・速度・半径・陣営を持ちます。", "データを追加するだけで、後ろの共通ループが自動的に動かしてくれます。"],
      ["異なる周期で二つの弾幕を呼ぶ", "38tickの円形弾と91tickの狙い扇を独立に発生させます。同じtickに重なれば両方出て、自然に複雑な模様になります。", "単純な周期の重ね合わせから複雑さが生まれるのが、弾幕作りの面白さです。"],
      ["全弾を後ろから動かす", "末尾からindexを減らし、ポインタで弾のx/yへ速度を足します。後で削除しても未処理indexがずれません。", "一発の規則をforで全弾へ広げるのがプログラムの力です。"],
      ["まず画面外かをremoveへ記録する", "Outsideの答えをboolへ入れます。この後の命中判定もremove=trueへ合流でき、削除場所を一か所にできます。", "『今すぐ消す』でなく『最後に消す印』を付ける設計です。"],
      ["自弾ならボスとの円判定を行う", "enemyでない弾だけボスを調べ、命中でHPを減らしremoveを立てます。HPが0以下ならclearedです。", "!b.enemyは、同じbullet型を陣営で使い分けます。"],
      ["敵弾なら無敵時間と自機判定を行う", "enemyで、無敵0で、円が重なる時だけlifeを減らします。無敵90tickを入れるため一度の接触で命が連続消費されません。", "条件を&&で重ねるほど、『本当に被弾する時』へ絞り込めます。"],
      ["印の付いた弾を一度だけ削除する", "画面外・ボス命中・自機命中のどれでもremoveならsliceから切り取ります。全ての道が一つの後始末へ集まります。", "更新と削除判断を分けると、条件を増やしても片付け方は変わりません。"],
    ]},
    en: { title: "Advance a bullet storm one shot at a time", lead: "Bullet hell combines spawn schedules with one rule per bullet. A backward slice pass handles movement, hits, and deletion.", blocks: [
      ["Close terminal states", "Clear or game over accepts restart only, freezing spawning and movement while Draw presents the result.", "The first return establishes every later assumption."],
      ["Advance the clock and invincibility", "frame schedules attacks; invincible spaces damage. It stops safely at zero.", "Different systems can reuse the same tick-counter idea."],
      ["Finalize player position first", "movePlayer runs before firing and aiming, so later rules use this tick's latest position.", "Choosing update order chooses game feel."],
      ["Append two player shots every seven ticks", "%7==0 is the firing clock; two struct literals carry position, velocity, radius, and team.", "Adding data is enough—the shared loop below will move it."],
      ["Call two patterns on different periods", "A 38-tick ring and 91-tick aimed fan may overlap, producing complexity from simple schedules.", "Layering simple periods is part of the fun of bullet patterns."],
      ["Move every bullet backward", "Count indices down and mutate x/y through a pointer, keeping deletion safe later.", "One rule expanded by for is the power of programming."],
      ["Mark off-screen bullets for removal", "Store Outside in a bool so later hit paths can join at remove=true and one cleanup point.", "Mark now, delete once at the end."],
      ["For player shots, test the boss", "Only !enemy bullets reduce boss HP, mark removal, and set cleared at zero.", "One bullet type can represent teams with a bool."],
      ["For enemy shots, test invincibility and player", "Team, zero invincibility, and circle hit must all pass; then life drops and recovery starts.", "Each && narrows the rule to a real hit."],
      ["Delete every marked bullet once", "Off-screen and either hit outcome all converge on one slice operation.", "New removal reasons do not need a new cleanup mechanism."],
    ]},
  },
];

function escapeHTML(value) {
  return String(value).replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;").replaceAll('"', "&quot;");
}

function updateBlocks(sourcePath) {
  const source = readFileSync(join(root, sourcePath), "utf8");
  const signature = "func (g *game) Update() error {";
  const start = source.indexOf(signature);
  if (start < 0) throw new Error(`${sourcePath}: Update signature not found`);
  let depth = 0;
  let end = -1;
  for (let i = start; i < source.length; i++) {
    if (source[i] === "{") depth++;
    if (source[i] === "}" && --depth === 0) { end = i + 1; break; }
  }
  if (end < 0) throw new Error(`${sourcePath}: Update closing brace not found`);
  return source.slice(start, end).split(/\n\s*\n/);
}

function render(lesson, lang) {
  const copy = lang === "ja" ? "この部分をコピー" : "Copy this part";
  const copied = lang === "ja" ? "コピーしました！" : "Copied!";
  const c = lesson[lang];
  const blocks = updateBlocks(lesson.source);
  if (blocks.length !== c.blocks.length) {
    throw new Error(`${lesson.source}: Update has ${blocks.length} blocks, but ${lang} lesson explains ${c.blocks.length}`);
  }
  const cards = blocks.map((code, index) => {
    const [title, body, note] = c.blocks[index];
    return `<article class="update-build-step"><header><span>${String(index + 1).padStart(2, "0")}</span><div><small>${lang === "ja" ? "PASTE NEXT / 次に貼る" : "PASTE NEXT"}</small><h3>${escapeHTML(title)}</h3></div></header><p>${escapeHTML(body)}</p><div class="update-go-note"><b>GO</b><span>${escapeHTML(note)}</span></div><div class="full-code update-code"><div class="full-code-head"><span>Update · block ${index + 1}/${blocks.length}</span><button type="button" class="full-code-copy" data-copy data-copied-label="${copied}">${copy}</button></div><pre><code>${escapeHTML(code)}</code></pre></div></article>`;
  }).join("\n");
  return `${markerStart}\n<section class="core-update-workshop" id="update-workshop"><div class="update-workshop-intro"><p class="eyebrow">${lang === "ja" ? "BUILD IT LINE BY LINE / UPDATE工房" : "BUILD IT LINE BY LINE / UPDATE WORKSHOP"}</p><h2>${escapeHTML(c.title)}</h2><p>${escapeHTML(c.lead)}</p><div class="update-workshop-route"><b>1</b><span>${lang === "ja" ? "上の「ぜんぶのコード」を main.go へ貼る" : "Paste the full source above into main.go"}</span><i>→</i><b>2</b><span>${lang === "ja" ? "そのUpdateを空にし、下の断片を番号順に貼る" : "Empty its Update, then paste the pieces below in order"}</span><i>→</i><b>3</b><span>${lang === "ja" ? "最後の断片まで貼って保存し、動かす" : "Save and run after the final piece"}</span></div><p class="update-workshop-tip">${lang === "ja" ? "途中は波かっこが閉じていない場合があり、まだ実行できなくて正常です。各カードでは「何を貼るか」だけでなく、「なぜこの順番か」「Goのどの書き方が役立つか」を一つずつ確認します。コードは完成ゲームの実物から取り出しているため、最後まで貼ると本物のUpdateと同じになります。" : "Some braces intentionally stay open between pieces, so intermediate code may not compile yet. Each card explains what to paste, why it belongs here, and which Go feature helps. The snippets are extracted from the running game; after the final paste, Update exactly matches the real one."}</p></div><div class="update-build-list">${cards}</div><div class="update-workshop-finish"><strong>${lang === "ja" ? "できあがり" : "COMPLETE"}</strong><p>${lang === "ja" ? "Updateは入力とゲーム状態を進める仕事だけを持ちます。Drawを別の絵へ交換しても、ここで組み立てた操作・時間・衝突・勝敗は変わりません。" : "Update owns only input and advancing game state. Replace Draw with another presentation and the controls, timing, collisions, and outcomes assembled here remain unchanged."}</p></div></section>\n${markerEnd}`;
}

for (const lesson of lessons) {
  for (const lang of ["ja", "en"]) {
    const path = join(root, "web", lang, "games", lesson.slug, "index.html");
    let html = readFileSync(path, "utf8");
    const existing = new RegExp(`${markerStart}[\\s\\S]*?${markerEnd}\\n?`, "g");
    html = html.replace(existing, "");
    const anchor = /\n\s*<div class="why-grid">/;
    if (!anchor.test(html)) throw new Error(`${path}: why-grid insertion anchor not found`);
    html = html.replace(anchor, `\n${render(lesson, lang)}\n\n    <div class="why-grid">`);
    writeFileSync(path, html);
  }
}

console.log(`Injected paste-in-order Update workshops into ${lessons.length * 2} core pages.`);
