#!/usr/bin/env node
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join, relative } from "node:path";

const root = new URL("../web", import.meta.url).pathname;
const marker = /<!-- BEGIN BEGINNER BRIDGE -->[\s\S]*?<!-- END BEGINNER BRIDGE -->\s*/g;

// Short, concrete bridges requested by classroom review.  These complement the
// deep dive: first define the unfamiliar word, then show one rule to trace.
const notes = {
  "setup":"go.modは『このゲームの名前と使う部品』を書いた目次ファイル。PowerShell/ターミナルは文字で命令する同じ役目の窓です。まずは空のゲーム画面を出すところまで進みます。|func (g *Game) Update() error { return nil }",
  "game-data":"データ駆動とは、村・台詞・敵をコードへ直接書かず、交換できるデータとして読む作り方。asset loaderはそのファイルをゲームへ渡す係、IDは部品同士を結ぶ名前札です。|village := loadMap(\"village.json\")",
  "ebi-depths":"巨大な世界は、世界座標・カメラ・訪問済み部屋・能力フラグを別々に覚えます。Updateで主人公と能力を進め、Drawではカメラの窓に入る部分だけを描きます。|screenX := worldX - cameraX",
  "timing-meter":"状態は『狙っている・判定を見ている・再挑戦待ち』のような今の場面。switchは状態ごとの処理をif/elseより読みやすく並べます。%は割り算の余りです。|phase := frame % 120",
  "catch-stars":"スライスは星を何個でも入れられる順番付きリスト。天次郎が星を受け止め、3個落とすとゲームオーバーです。|stars = append(stars, newStar())",
  "flappy":"画面は上がY=0で、下へ行くほどYが増えます。加速度は『Updateの刻みごとに速さがどれだけ変わるか』。羽ばたき時の負の速度は上向き、正の重力が少しずつ下向きへ戻します。|vy = -7.4; vy += gravity; y += vy",
  "pong":"ベクトルは横の速さvxと縦の速さvyを組にしたもの。符号反転は+を-へ変え、壁へ向かう速さを反対向きにします。|vx = -vx",
  "breakout":"[]は順番付きの箱、[][]は行と列の方眼紙。全部のbrickが消えたらwinがclear状態を立て、次のDrawで結果画面を出します。|if remaining == 0 { state = Win }",
  "snake":"スライスは長さを変えられるリスト。ヘビの体は通ったマス座標の順番で、appendで頭を足し、最後を外すと一歩進みます。|body = append([]Cell{head}, body[:len(body)-1]...)",
  "space-shooter":"自機・弾・敵を別リストにすると、それぞれ同じUpdateを繰り返せます。全弾×全敵を調べる総当たりは、教材の小さな数なら十分です。|for _, b := range bullets { for _, e := range enemies { check(b,e) } }",
  "sokoban":"論理座標は『何列・何行』、描画座標は画面のピクセル。天次郎が1マス動く間はmoving=trueにして入力を待たせ、Drawだけ途中位置を描きます。|if moving { progress++; return nil }",
  "platformer":"めり込みは物体が壁の中へ重なること。Xを動かして直し、次にYを動かして直すと角でも滑れます。0.09はカメラが差の9%ずつ追う調整値です。|cameraX += (playerX-cameraX) * .09",
  "dungeon":"patrol（巡回）は決まった道の見回り、chase（追跡）は天次郎を追う状態。165pxなどは気付く距離を遊びながら調整した値です。|if distance < 165 { state = Chase }",
  "bullet-hell":"cos/sinは角度から円周上の横・縦位置を得る道具。ラジアンは角度単位で2πが一周、tick%周期は一定tickごとに同じ模様を繰り返します。|angle := 2*math.Pi*float64(i)/count",
  "moving-platforms":"相対移動は『床が前のコマから何px動いたか』という差を天次郎にも足す考え方。動かす前のoldXを記録し、dx=新しいx-oldXを求めます。|oldX:=platform.x; platform.x+=platform.vx; player.x+=platform.x-oldX",
  "patrol-enemies":"巡回は敵が二つの端点を往復する見回り。踏み判定は、天次郎の前の足位置が敵より上で、今の足位置が敵の上端を越えたかを調べ、上から当たった時だけ敵を倒します。|stomp := oldBottom <= enemy.y && newBottom >= enemy.y",
  "powerup-adventure":"状態管理は今の強化をGameへ保存すること。bigは『大きい姿か』を覚えるtrue/falseの印です。取得矩形が重なったらitem.kindに応じてspeed・jump・HP・bigを変えます。|if item.kind==Grow { g.big=true }",
  "scrolling-stage":"ワールド座標は広いステージ内の位置、画面座標は今見えている窓の位置。screenX=worldX-cameraXで変換します。カリングは画面外の物をDrawしない工夫です。|screenX := worldX - cameraX",
  "tiny-platformer":"接地判定は天次郎の足が床へ触れたか調べること。触れたら落下速度を0にし、ジャンプ入力で上向き速度を入れます。赤い四角は最初に動きを確認する仮のプレイヤーです。|if onFloor && jump { vy=-jumpPower }",
  "arena-dodge":"正規化は斜め移動だけ速くならないよう方向の長さを1へそろえること。左右dxと上下dyの長さを求め、0でなければ両方をその長さで割ります。|length:=math.Hypot(dx,dy); dx/=length; dy/=length",
  "auto-turret":"最近傍探索は全敵の距離を比べて一番近い敵を選ぶこと。クールダウンは次に撃てるまでの待ち時間で、Updateごとに1減らし0で発射します。|target:=nearestEnemy(); if cooldown==0 { fire(target) }",
  "experience-draft":"XPは敵を倒した経験値、レベルアップはXPが境目を越えること、ドラフト選択はランダム候補から一つ強化を選ぶことです。条件を満たしたら時間を止めて3択状態へ移ります。|if xp>=nextXP { state=ChooseUpgrade }",
  "survival-run":"ウェーブは時間で区切った敵の波、ボスは区切りの最後に出る強敵、難易度曲線は時間とともに数・速さ・HPを少しずつ増やす設計です。目標時間まで生き残ればクリアです。|wave := elapsed / waveSeconds",
  "swarm":"全敵を一度ずつUpdateするので敵数に比例して仕事が増えます。画面外描画を省き、上限を決めて60fpsを守ります。|for i := range enemies { enemies[i].Update() }",
  "weapon-evolution":"データ駆動は強さをifの山でなく武器データ表に書く方法。経験値ゲージが満ちるたびlevelを増やし、2で3way、3で連射など次のデータへ交換します。合成は二つの条件を満たして上位武器にすることです。|weapon = evolutions[level]",
  "command-battle":"状態機械は『入力待ち→味方行動→敵行動→結果』の今の場面を一つのstateで覚える仕組み。A/B/Cまたは画面ボタンをAttack/Guard/Skillへ翻訳し、ターンを順に進めます。|switch state { case Choose: readCommand(); case Enemy: enemyTurn() }",
  "dialogue-flags":"フラグは出来事を覚える名前付きのtrue/false。文章送りはボタンで次の台詞番号へ進むこと、イベントフラグは選択後も残り、後の台詞をifで分けます。|if choice==Help { g.helped=true }",
  "ebi-quest":"仲間は一緒に戦う人物データ、クエストは目的と進み具合の組、セーブはその状態を後で読める形で残すこと。クエスト番号で今の目的を選びます。|quest := quests[questID]",
  "inventory-shop":"所持品は持っている道具一覧、装備は今使って能力へ足す一つ、売買はgoldとitemを交換する処理。equippedIDと同じなら装備済みです。|if item.id==equippedID { label=\"装備中\" }",
  "stats-status":"ダメージは基本的に攻撃-防御で、最低1に丸めます。一時効果は数ターンだけ攻撃や防御を増減する値。structはHP/ATK/DEFと残りターンを一組にする箱です。|damage:=max(1, attacker.ATK-defender.DEF)",
  "village-walk":"タイル衝突は移動先の地図文字が壁か調べること、向きは天次郎が上下左右のどちらを見ているか。rows[y][x]を読み、通路なら座標と向きを更新します。|if rows[ny][nx]!='#' { x,y,dir=nx,ny,inputDir }",
  "world-encounters":"シーン変数は『村・フィールド・戦闘』の今の画面番号。敵テーブルは場所ごとに出せる敵の候補一覧、シーン遷移は番号を替えて別のUpdate/Drawへ渡すことです。|scene=Battle; enemy=enemyTable[area].Pick()",
  "tap-counter":"クリック入力は押した瞬間を読むこと、カウンターは回数を覚えるcount。inpututil.IsMouseButtonJustPressedは左ボタンが今のコマで押された瞬間だけtrueを返します。|if inpututil.IsMouseButtonJustPressed(0) { g.count++ }",
  "first-shop":"購入条件sweets>=costは『持っているお菓子が価格以上』という比較。買えたらお菓子を減らし、ownedを増やし、base+owned×stepで次の価格を上げます。|if sweets>=cost { sweets-=cost; owned++ }",
  "growth-curves":"成長曲線は回数と必要価格の関係。price=base×1.15^ownedなら買うたび約115%。1.00MのMはmillion（100万）の略で、表示だけ短くし本当の数は変えません。|price := base * math.Pow(1.15, float64(owned))",
  "idle-factory":"デルタタイムdtは前のコマから何秒たったか、CPSはCookies Per Secondのような一秒あたり生産量。生産量×dtを足すと画面速度が変わっても同じ量になります。|stock += cps * dt",
  "offline-bakery":"セーブは状態を保存、ロードは復元、オフライン進行は閉じていた時間ぶんを再開時に計算すること。localStorageはこのブラウザ内だけに文字を保存する棚です。|gain := time.Since(savedAt).Seconds()*perSecond",
  "move-battle":"一ターンをPlayerMove→PlayerAttack→EnemyMove→EnemyAttackのphaseで進め、各phaseが終わるまで次を動かしません。|phase = PlayerMove",
  "ebi-monsters":"完成形はSpecies・Stats・Moves・Partyの四箱から読みます。まず一匹のspeciesを選び、能力、技、仲間順にたどります。|monster := Monster{Species: speciesID, Stats: baseStats[speciesID]} ",
  "inventory-crafting":"ドラッグ中はselectedItemを覚え、離したマス番号へ移します。レシピは配置後のID列と必要ID列を比べます。|if released { slots[target] = selectedItem }",
  "tools-light":"Tool structはNameとHardnessを一組にする箱。tool.Hardness>=block.Requiredなら壊せます。|if tool.Hardness >= block.Required { breakBlock() }",
  "creature-world":"noiseは近い場所ほど近い値を返す乱数風の関数。値の範囲を草・砂・森へ割り当て、seedで同じ世界を再現します。|biome := biomeAt(seed, x, y)",
  "ebi-craft":"完成形は入力→座標変換→地図更新→持ち物更新→Drawの順。各STEPの関数をこの順に呼び、巨大な一関数にしません。|handleInput(); updateWorld(); updateInventory()",
  "pellet-map":"tile値0=壁、1=通路、2=えさと決め、天次郎のマスが2なら1へ変えて得点します。|if maze[y][x] == Pellet { maze[y][x]=Floor; score++ }",
  "tile-runner":"一マスずつのリアルタイム移動で、天次郎→敵の順に更新します。移動先map[ny][nx]が壁でなければ座標を採用します。|if map[ny][nx] != Wall { x,y = nx,ny }",
  "ebi-maze":"完成形の地図は固定タイルデータから読み、移動・えさ・AIを別関数にします。迷路生成アルゴリズムはこの教材の必須部品ではありません。|updatePlayer(); collectPellet(); updateEnemies()",
  "escape-ai":"未来の爆風マスを危険印にし、BFSは近いマスから順に安全マスを探します。queueへ未訪問の上下左右を足します。|queue = append(queue, neighbors(current)...)",
  "ebi-bomber":"完成形も設置→timer更新→爆風計算→壁/敵判定→出口判定の順。各STEPの関数をそのまま一列に並べます。|placeBomb(); tickBombs(); resolveBlasts()",
  "ebi-merge":"予告は次に落とす実、得点は合体の成果、ゲームオーバーは置き場所がなくなった状態。落下中は左右を考える時間を残し、同じtierの実が触れた時だけ合体します。合体済みへremoved印を付けて二重処理を防ぎます。|if a.tier==b.tier { merge(a,b); score+=value }",
  "ebi-blocks":"形[][]、回転候補[][][]、現在位置を別に持ちます。[]が一段増えるごとに『形の一覧』が一層増えると読めます。|rotatedShape := rotations[piece][rotation]",
  "ally-effects":"前STEPのParticle型とspawn関数を同じpackageから呼び、仲間の種類だけ色・速さ・個数の引数で変えます。|spawnEffect(hitX, hitY, ally.effect)",
  "circle-battle":"全敵を一度ずつ調べ、中心距離²が半径合計²未満なら命中。敵数が小さい教材ではこのO(n)走査が最も読みやすい方法です。|for i := range enemies { checkCircle(projectile, enemies[i]) }",
  "ebi-strike":"完成形はaim→launch→move→collide→fxの五関数。どのSTEPの部品かが呼び出し順で分かります。|aim(); launch(); move(); collide(); updateFX()",
  "falling-pair":"入力中は列だけ選び、Updateの落下tickで二個のrowを同時に+1。Drawは親と相方の二座標を同じ進み具合で描きます。|pair.y += fallSpeed",
  "ebi-chain":"完成形はinput→fall→lock→findGroup→clear→gravityをphaseで一つずつ進めます。最初はこの一列を復習地図として読みます。|updateInput(); updateFall(); resolveBoard()",
  "play-a-card":"カードデータは名前・必要エナジー・効果を一組にした設計図。効果解決はカードの指示どおりHPや防御を変える処理です。同じresolve関数へ違うカードを渡すと、同じ手順で別効果を実行できます。|resolve(card)",
  "turn-energy":"ターンフェーズは『自分が選ぶ・敵が動く・次の手札を配る』という場面分け。行動資源は一ターンに使えるエナジーで、カードを使うたび減ります。その場面のUpdateだけ動かします。|if phase == PlayerTurn { updatePlayerTurn() }",
  "deck-cycle":"山札はこれから引くリスト、手札は今使えるリスト、捨て札は使用済みリスト。カードを山札→手札→捨て札へ移すことで、データの居場所が今の状態を表します。|hand = append(hand, drawPile[0])",
  "intent-status":"予告行動は敵が次に何をするかの表示、持続期間は毒や弱体があと何ターン残るか。Updateで残りターンを減らし、Drawで予告の絵と数字を見せます。|status.turns--; drawIntent(enemy.intent)",
  "card-rewards":"選択報酬は戦闘後に一枚だけ選ぶごほうび、レア度は出にくさと強さの目安、デッキ編集は使うカード一覧を変えること。選んだ一枚だけappendします。|deck = append(deck, chosenCard)",
  "branching-map":"グラフは場所（ノード）と道（つながり）の組。nextは今いる場所から次に選べる場所番号のリストです。番号を選ぶとcurrentをその番号へ変えます。|current = nodes[current].next[choice]",
  "card-motion":"PLAYABLEはその場で遊べる教材という印。デッキ構築カードゲームは、戦いながらカードを選び自分の山札を育てるゲームです。このSTEPはカードの移動→命中→戻りを中間フレームで見せます。|phase = CardTravel",
  "run-director":"PLAYABLEはその場で遊べる教材という印。run directorは次の敵・休憩・報酬・背景を階層データから選び、一回の冒険を進める係。Ebi Ascentは完成ゲーム名です。|encounter := floors[floor]",
  "ebi-ascent":"phaseは『戦闘中・報酬選択・道選択・休憩』など現在の場面。状態をつなぐとは、一つを終えたら次のphase番号へ替えてゲーム全体を進めることです。分岐マップは次の行き先を選べる旅の地図です。|switch phase { case Battle: updateBattle() }",
  "run-animation":"PLAYABLEはその場で遊べる教材という印。横スクロールアクションは天次郎を横へ進めて跳ぶゲーム。走りアニメは足の絵を時間で切り替え、速さに応じて再生間隔を変えます。|frame := tick / frameDuration % frameCount",
  "stage-data":"ステージデータは床・敵・ゴールの位置をまとめた設計図。横スクロールアクションはそのデータの一部分をカメラで見せます。海老・天次郎は同じUpdateのまま別ステージを走れます。|stage := stages[stageIndex]",
  "game-feel":"PLAYABLEはその場で遊べる教材という印。game feelは、同じルールでも揺れ・粒・停止のタイミングで操作を気持ちよく伝える工夫。Ebi Adventureは完成ゲーム名です。|hitStop=4; cameraShake=6",
  "impact-feedback":"PLAYABLEはその場で遊べる教材という印。大量敵サバイバルは囲まれる敵を避けながら自動攻撃で生き残るゲーム。impact feedbackは命中の点滅・粒・小さな画面揺れです。|enemy.flash=6; shake=4",
  "wave-director":"PLAYABLEはその場で遊べる教材という印。wave directorは時間に応じて敵の種類・数・ボスを選ぶ係。難易度曲線は序盤から終盤へ少しずつ難しくする並びです。|wave := waves[elapsed/30]",
  "vfx-stamp":"GeoM は絵の置き場所を覚える道具です。Translate(x,y) の x は右、y は下へ動かす量。& は設定箱そのものを渡す印です。|op.GeoM.Translate(x, y)",
  "vfx-transform":"GeoM.Rotateは絵を回す設定、GeoM.Scaleは横・縦の大きさを何倍にするかの設定。w/hはwidth（横幅）/height（高さ）です。-w/2,-h/2で絵の中心を座標(0,0)へ合わせてから回します。スタンプと移動で使ったTranslateへ、回転と拡大縮小を足す段階です。|op.GeoM.Rotate(angle); op.GeoM.Scale(sx, sy)",
  "vfx-tint":"tint は元絵へ色を重ねること。ColorScale は赤・緑・青・透明度の順で、1はそのまま、0.4は40%残す意味です。|op.ColorScale.Scale(1, .4, .4, 1)",
  "vfx-alpha":"alpha（アルファ）は透明度の数字で、0なら見えず、1なら元の濃さ。ScaleAlpha(.5)は透明度を半分にします。GeoMは絵の置き場所を変える設定で、古い位置を薄くしてから現在位置を描くと残像になります。|op.ColorScale.ScaleAlpha(.5)",
  "vfx-additive":"加算合成は、すでにある光へ次の光の明るさを足す描き方。ebiten.BlendLighterは『重なった色を明るい方へ足す』設定です。赤い光と緑の光は重なると黄色へ近づきます。前の透明度・重ね順・残像の描き方へ光らせ方を追加します。|op.Blend = ebiten.BlendLighter",
  "vfx-walk":"SubImage（サブイメージ）は、大きな画像から四角い一部分を切り出す命令。タイマーは何コマ進んだかを数える整数（小数でない数）です。type は、その数や絵をまとめる新しい箱の種類を決める書き方。col は横、row は縦の何番目かで、col×コマ幅が左端になります。Updateがコマ番号を変え、Drawが切り出した絵を表示します。加算合成は光を足す描き方、パーティクルは短時間だけ動く小さな粒です。|frame := sheet.SubImage(frameRect).(*ebiten.Image)",
  "vfx-particles":"particleは一粒分の位置・速さ・残り時間をまとめた型。[]particleはその粒を順番に並べるリストです。Updateで全粒のlifeを減らし、0より大きい粒だけ残し、Drawで現在位置へ描きます。前のコマ送りは一枚ずつ絵を替えましたが、ここでは粒を何個も同時に進めます。|if p.life > 0 { next = append(next, p) }",
  "vfx-spells":"エフェクトの合成は、別々に作った光・魔法本体・火花を順番に重ねて一つに見せること。ScaleAlphaは透明度を何倍にするか、ColorScaleは赤・緑・青・透明度を何倍にするかの設定です。整数型は小数を持たない数の種類、typeは新しいデータの種類を定義するGoの言葉。Updateが粒の数字を進め、Drawが背景の光→本体→粒の順に描きます。パーティクルは生まれて動いて消える小さな粒です。|drawGlow(); drawSpell(); drawParticles()",
  "vfx-magic-fire":"レイヤーは奥から手前へ重ねる絵の層。チャージ炎は、押しているフレーム数を力として、奥の光→炎本体→手前の火花の順に描きます。しきい値は段階が変わる境目で、60フレームなら約1秒です。|drawBackGlow(); drawFire(); drawFrontSparks()",
  "vfx-magic-ice":"凍結を霜・氷殻・停止・亀裂・粉砕に分けます。水色の殻を一拍止めてから割るので、青い爆発ではなく「動きを封じた」と読めます。|if phase == Freeze { drawShell() } else { shatter() }",
  "vfx-magic-thunder":"枝分かれ稲妻は、一本の線の途中から短い線を左右へ伸ばした形。再帰は線を分ける処理が自分自身をもう一度呼ぶ書き方です。角度は見た目の調整値なので、まず±30度から比べます。|branch(a, b, depth-1)",
  "vfx-magic-light":"神聖さは明るさだけでなく、上から降りる光柱、動かない輪光、左右対称の線、ゆっくり落ちる金色の粒で作ります。ランダムな爆発を秩序ある儀式へ変える例です。|DrawPillar(); DrawHalo(); fallBlessingMotes()",
  "vfx-magic-dark":"渦は中心の周りを回る点の列、吸い込みは中心までの距離rを毎コマ減らす動き。角度を増やしながらrを減らすと、らせん状に中心へ入ります。強さ2.0は一度に近づく量です。|r -= 2; angle += .12",
  "vfx-fx-tap":"押した瞬間だけ粒を出すにはJustPressedを使います。isActivatedは同じ的から二度ごほうびを出さないための印です。|if justPressed && !isActivated { isActivated=true; fx.Burst() }",
  "vfx-fx-meter":"『状態ごとに違う演出』とは、判定結果をPerfect・OK・Missの三状態として覚え、光の色・粒の数・揺れを変えること。Updateが針と判定の数字を決め、Drawがその状態に合う絵を描きます。LEVEL 02と同じ差を使い、中心から±4ならPerfect、±12ならOK、それ以外はMissです。|result := judge(math.Abs(marker-center))",
  "vfx-fx-catch":"エンティティは星や天次郎のようにゲームのルールへ参加する物体、FXは勝敗を変えない一瞬の見た目。混ぜないとはstarsとparticlesを別リストで持つことです。Updateは星の落下と捕獲を決め、Drawは星と粒を重ねます。appendはリスト末尾へ新しい一個を足します。|particles = append(particles, burst...)",
  "vfx-fx-flight":"重力・関門・得点は勝敗を変えるplay、羽・通過リング・光は見た目だけのfx。分けると演出を消しても同じゲームルールを保てます。|updatePlay(); fx.Update()",
  "vfx-fx-pong":"衝突はボールとパドルの四角が重なること。トレイルはボールが通った場所へ薄い絵を短く残す残像です。Updateの衝突判定で速度を反転して火花を生み、Drawで古い残像→ボールの順に描きます。|if hitPaddle { vx=-vx; fx.Spark(x,y) }",
  "vfx-fx-breakout":"エンティティはルールに参加するブロック、FXは破片など見た目だけの演出。『破壊=削除+FX』は、Updateで命中ブロックへalive=falseを付け、同じ場所に破片を生む二つの仕事です。Drawは生きたブロックと破片だけを描きます。|if hit { brick.alive=false; fx.Shards(brick.x,brick.y) }",
  "vfx-fx-snake":"セルは方眼紙の一マス、体セルはヘビの体があるマス。キラキラは食べた場所で短時間だけ動くFXです。[]Cellは体のマス一覧。Updateで頭を足して尾を外し、Drawで食べ物→体→粒を描きます。|drawFood(); drawBody(); fx.Draw()",
  "vfx-fx-shooter":"三分割とは、弾bullets・敵enemies・爆発FXを別リストで管理すること。弾は命中判定、敵はHPと移動、FXは見た目だけを担当します。Updateで三つを順に進め、Drawで敵→弾→爆発の順に重ねます。|updateBullets(); updateEnemies(); fx.Update()",
  "vfx-fx-sokoban":"タイル地図は壁・床・ゴールを番号で並べた方眼紙。押し演出は箱が一マス動いた瞬間だけ出す塵です。Updateでmap[y][x]を書き換え、Drawで地図→箱→塵→ゴールの輝きを描きます。|if pushed { fx.Dust(boxX, boxY) }",
  "vfx-fx-platform":"接地は天次郎の足が床へ触れた状態、イベントは状態が『空中→接地』へ変わった一瞬。塵はその一瞬に生む粒です。Updateで接地印を比べ、Drawで天次郎と塵を重ねます。|if !wasGrounded && grounded { fx.Dust(x,y) }",
  "vfx-fx-dungeon":"ヒット反応は攻撃が当たったことを伝える点滅・のけぞり・粒のFX。HPはゲームの本当の数なのでUpdateで減らし、flashは見た目用タイマーとして分けます。Drawはflash中だけ敵を明るく描きます。|enemy.hp -= damage; enemy.flash = 8",
  "vfx-fx-bullethell":"弾幕は多数の弾で避け道を作る仕組み、ボムは画面内の弾を消す救済技、FXはその爆発を見せる演出です。Updateで弾リストを空にして勝敗へ反映し、Drawで爆発の光と粒を描きます。frame%120は120コマごとの周期です。|bullets = bullets[:0]; fx.Bomb(playerX,playerY)",
  "cascade":"消す→落とす→もう一度探すをphaseで分け、揃いがなくなるまで往復すると連鎖になります。|phase = Find; if matched { phase = Clear }",
  "clear-and-fall":"論理盤面はピースが最終的にいる行列、見た目のyは落下途中の画面位置。Updateでprogressを0→1へ進め、Drawで前のyと次のyの間をprogressの割合だけ補間します。落下完了後に盤面を確定します。|drawY := oldY + (newY-oldY)*progress",
  "find-matches":"走査は方眼紙を端から順に見ること。横はy固定でxを、縦はx固定でyを進め、揃った座標をmarkedという重複しない集合へ入れます。|marked[Cell{x,y}] = true",
  "grid-swap":"2次元グリッドはboard[行][列]で読む方眼紙、隣接判定は列差+行差が1かを調べること。隣同士だけa,b=b,aで交換し、3個並びがなければ戻します。|adjacent := abs(ax-bx)+abs(ay-by)==1",
  "special-pieces":"形の判定は揃った座標が直線・L字・T字か読むこと、範囲効果は特殊ピースが消すマス一覧。①形からkindを決め、②kindごとに消す座標集合を返します。|cells := affectedCells(piece.kind, x, y)",
  "ebi-match":"手数は交換できる残り回数、目標は集める色と個数、ステージデータは盤面・手数・目標をまとめた設計図。ルールは共通のUpdateに置き、数値だけデータから読みます。|stage := stages[stageIndex]",
  "box-viewer":"矩形は位置x,yと幅w,高さhを持つ四角。押し合い判定は体同士が通り抜けない箱、食らい判定は攻撃を受ける箱、攻撃判定は技を出す間だけ現れる箱です。|hit := attackBox.Overlaps(hurtBox)",
  "combo-cancel":"入力バッファは少し早く押した次の技を予約する箱、キャンセルウィンドウは今の技から次へ移れるフレーム範囲。8〜15の間だけ予約を実行します。|if frame>=8 && frame<=15 { acceptBufferedMove() }",
  "command-fighter":"入力履歴は押した方向・ボタン・時刻を古い順に残すリスト、コマンド認識は末尾の並びが技の決められた順番と一致するか調べること。20フレームより古い入力を捨ててから比較します。|inputs = append(inputs, Input{key, frame})",
  "ebi-fighters":"ラウンドはどちらかの体力が0になるまでの一勝負、体力はあと何回攻撃に耐えられるかの数、対戦ルールは勝利条件と時間制限。ラウンド勝利数を先に2へした側が試合に勝ちます。|if hp<=0 { roundWins[winner]++ }",
  "frame-attack":"60fpsで4フレームは約0.067秒。攻撃矩形はactiveの間だけ出し、その前後は当たりません。|active := frame >= 6 && frame < 10",
  "guard-throw":"まず技の矩形どうしが重なったかを確かめます。重なった接触の中で、投げはガードに勝ち、ガードは打撃に勝ち、打撃は投げの準備に勝つ優先順位を使います。離れて矩形が重ならなければ、どの技も空振りです。|if boxesOverlap(attackBox, targetBox) { resolvePriority() }",
  "hit-reaction":"攻撃矩形と食らい矩形が重なった一瞬にHPを減らし、同じ出来事からのけぞり演出を始めます。|if attack.Overlaps(hurt) { hp--; state=Hurt }",
  "circle-collision":"中心距離は二円の中心どうしの離れ具合、接触法線は一方の中心からもう一方へ向く長さ1の矢印。中心距離が半径二つの合計より小さければ重なっています。|dx*dx+dy*dy < (ra+rb)*(ra+rb)",
  "falling-circles":"固定時間ステップは毎Updateを同じ短い時間として計算する方法。同じ重力g=.5を毎コマ速度へ足し、円の下端y+rが床を越えたら床-rへ戻します。|vy += .5; y += vy",
  "merge-rule":"同じtierの円が触れたら二つを消し、中心へtier+1を一つ作ります。|if a.tier == b.tier { merge(a, b) }",
  "stable-simulation":"反復解法は小さなめり込み修正を7回など複数回繰り返す方法、休眠はほぼ止まった円の計算を休む印、空間分割は近い円だけ調べる区画分け。回数は速さと安定を比べて決めます。|for pass:=0; pass<7; pass++ { solveContacts() }",
  "stacking-bounce":"めり込み修正は重なった円を接触線に沿って離すこと、インパルスは衝突の瞬間だけ速度へ加える押し返し。反発係数.8なら速さの80%で反対へ戻ります。|velocity += collisionNormal * impulse",
  "bag-hold-ghost":"公平な抽選7-bagは7種類を一つずつ袋へ入れ順番だけ混ぜる方法。ホールドは現在の形と保管形を状態交換し、ゴーストは衝突する直前まで仮に下げて着地点を予測します。|held,current = current,held",
  "tetromino-shapes":"相対座標は形の中心から何マスずれたかを表す座標。四つのずれを形データにすると、基準位置へ足すだけで4マスの形を置けます。|boardX := pieceX + cell.dx",
  "rotation-kicks":"座標回転は(dx,dy)を(-dy,dx)へ替えること。壁蹴り候補は回転後に重なった時のずらし案で、0、左1、右1、左2、右2を順に試します。|rotatedDX,rotatedDY := -dy,dx",
  "lock-lines":"盤面への合成は落下中の四マスをboardへ固定すること、行の圧縮は満杯行より上を一段下へコピーして空きを詰めることです。|lockPiece(); copy(board[1:y+1],board[0:y])",
  "falling-cell":"一定間隔更新は毎コマ落とさずtimerで待った回数を数える方法。timerがfallIntervalへ達した時だけ行を+1し、床なら固定します。|timer++; if timer>=fallInterval { row++; timer=0 }",
  "ebi-blocks":"状態の流れはSpawn→Fall→Lock→Clear→Spawnという小さな仕事の列。バッグは7種を公平に出す袋、ホールドは一個保管、ゴーストは着地点予測です。|switch phase { case Fall: updateFall(); case Clear: clearLines() }",
  "drag-launch":"押し始めから現在までの差dx,dyを取り、離したら反対向きの-dx,-dyへ倍率を掛けて速度にします。|vx, vy = -dx*power, -dy*power",
  "friction-stop":".95は現実の定数でなく操作感用。毎回95%へ減らし、十分小さければ0へ丸めます。|vx *= .95; if math.Abs(vx)<.05 { vx=0 }",
  "wall-bounce":"左右壁はvx、床天井はvyだけ符号反転し、位置も壁の内側へ戻して二重衝突を防ぎます。|vx = -vx; x = wallX-radius",
  "color-groups":"flood fillは一マスから上下左右の同色へ広がる探索。訪問済み印で同じマスを二度数えません。|visit(x,y); visit(x+1,y); visit(x,y+1)",
  "pair-rotation":"相方の差(dx,dy)を90度回すと(-dy,dx)。新しい列と行が盤内で空の時だけ採用します。|newDX, newDY := -dy, dx",
  "chain-score":"落下後にもう一度消えたらchainを増やし、後の連鎖ほど倍率を上げて先読みへごほうびを出します。|score += cleared * 100 * chain",
  "clear-gravity":"各列を独立に下から詰めます。下→上なら、まだ読んでいないピースを上書きしません。|for x:=0; x<width; x++ { compactColumn(x) }",
  "capture":"60%は正解でなく調整値。HPが少ないほど成功率を上げ、0〜99の乱数が確率未満なら捕獲です。|chance := 25 + (maxHP-hp)*60/maxHP",
  "growth-evolution":"経験値追加後、前level<5かつ新level>=5なら一度だけ進化を始めます。|if oldLevel < 5 && level >= 5 { evolve() }",
  "party-switch":"選んだ番号selectedと先頭0番を同時交換します。|party[0], party[selected] = party[selected], party[0]",
  "species-data":"mapは『名前札→データ』の表。species[\"coral\"]でcoralのStats一組を取り出します。|base := species[monster.species]",
  "type-matchup":"火→木→水→火を表で引き、有利1.5倍・不利.75倍にします。倍率は遊びやすさで変える値です。|damage *= matchup[attackType][defendType]",
  "chunk-world":"チャンクは広い地図の小区画。x/256が区画番号、x%256が区画内位置です。|chunkX, localX := x/256, x%256",
  "terrain-generation":"0=草、1=砂、2=水の対応表を先に決めます。同じseedなら同じ教材世界を再現できます。|tile := rng.Intn(3)",
  "place-break":"画面座標をTileSizeで割ると列と行。入力を『壊す』『置く』へ変換してから地図を更新します。|tx, ty := mouseX/TileSize, mouseY/TileSize",
  "buffered-turn":"tickごとに入力を見てwantedDirへ予約し、次の交差点で通れる時に採用します。|wantedDir = pressedDir",
  "junction-ai":"壁と逆戻りを候補から除き、プレイヤーへ近づく距離が小さい方向を選びます。|score := distance(next, player)",
  "patrol-chase":"patrolは見回り、chaseは追跡。距離100未満で追い、見失って時間が過ぎたら見回りへ戻ります。|if distance < 100 { state=Chase }",
  "cross-blast":"上下左右を一マスずつ射程まで進み、画面外か固い壁で止めます。壊れる壁は含めてから止めます。|for d:=1; d<=power; d++ { if blocked { break } }",
  "chain-explosion":"別爆弾へ届いたらタイマーを0にします。処理前にexploded=trueとしてA→B→Aの二重爆発を防ぎます。|if b.exploded { return }; b.exploded=true",
  "timed-bomb":"Updateごとにtimer--し、0なら爆発へ進みます。待ち時間でゲーム全体を止めません。|bomb.timer--; if bomb.timer <= 0 { explode() }",
  "breakable-walls":"爆風座標ごとにmap[y][x]を調べ、壊れる壁なら床へ変更。20%は0〜.999の乱数が.2未満かです。|if tile == Breakable { map[y][x] = Floor }",
};

const trackParts = {
  "visual-effects":"絵を置く→変形→色→重ね方→アニメ→粒", platformer:"動かす→床で止める→跳ぶ→カメラ", survivors:"移動→敵→自動攻撃→経験値→強化", clicker:"押して増やす→買う→自動生産→保存", rpg:"村→会話→戦闘→持ち物", fighting:"攻撃矩形→フレーム→反応→防御→連続技", "merge-physics":"落下→衝突修正→反発→合体", deckbuilder:"引く→コスト→効果→敵予告", match3:"交換→検索→消去→落下→再検索", "falling-blocks":"一マス→形→回転→固定→行消去", slingshot:"引く→速度→反射→命中", sandbox:"座標変換→壊す/置く→持ち物→生成→区画", "monster-collection":"種族→相性→戦闘→捕獲→成長", "falling-pairs":"二個組→回転→同色探索→落下→連鎖", "maze-chase":"マス移動→曲がり予約→餌→分岐AI→追跡", "bomb-maze":"タイマー→十字爆風→壁→連鎖→逃走", tactics:"移動範囲→地形コスト→射程→敵予告", "active-rpg":"素早さ→ゲージ→READY→コマンド→演出", "visual-novel":"会話→立ち絵→選択肢→フラグ→結末", racing:"加速→旋回→ゲート→ライバル→コース", metroidvania:"カメラ→部屋→地図→能力ゲート→戻り道"
};

const generic = [
  [/func \([^)]+\)/,"func (g *Game)の(g *Game)は『このGameの仕事』という名札（レシーバー）。g.hpのように、そのゲームが覚えた値を使えます。"],
  [/append|スライス/,"スライスは個数を増減できる順番付きリスト。append(list,item)は末尾へitemを一つ足します。"],
  [/\[\]\[\]|二次元配列/,"[][]は行の中に列がある方眼紙。まず[y]で行、次に[x]で列を選びます。"],
  [/struct|構造体/,"structは位置・速さ・HPなど関係する値を一つの箱へまとめる書き方です。"],
  [/Update|Draw/,"Updateは入力とルールからゲームの状態を進めます。Drawは呼ばれた時点の状態を読むだけで、進行を変えません。二つの回数や順番を同じものとは考えません。"],
];

function walk(dir){return readdirSync(dir).flatMap(n=>{const p=join(dir,n);return statSync(p).isDirectory()?walk(p):[p]})}
for(const path of walk(root).filter(p=>p.endsWith(".html"))){
  let html=readFileSync(path,"utf8").replace(marker,"");
  const route="/"+relative(root,path).replace(/index\.html$/,"").replaceAll("\\","/");
  if(!html.includes("<main")||!/(\/games\/|\/tracks\/|\/guides\/|^\/(ja|en)\/$)/.test(route))continue;
  const bits=route.split("/").filter(Boolean), lang=bits[0], slug=bits.length===1?"home":bits.at(-1);
  // Home must lead with an obvious choice, and LEVEL 01 must lead with play.
  // Their own pages introduce terminology only after the reader has context.
  if(slug === "home" || slug === "tap-target") {
    writeFileSync(path, html);
    continue;
  }
  const lookupSlug=slug.startsWith("fx-")?`vfx-${slug}`:slug==="three-punch"?"frame-attack":slug==="hit-reaction-dojyo"?"hit-reaction":slug;
  const isHub=bits[1]==="tracks"&&bits.length===3;
  let pair=notes[lookupSlug]?["ここを先に見よう",...notes[lookupSlug].split("|")]:undefined;
  if(isHub&&trackParts[slug]) pair=[`このコースの組み立て図`,`${trackParts[slug]}の順に、海老・天次郎（えび・てんじろう）と一部品ずつ作ります。`,`func (g *Game) Update() error { return nil }`];
  if(isHub&&trackParts[slug]){
    // “このページの1本” is the shared skeleton, not the finished game.
    pair[1] += ` このページの1本はコース全体で共通するgame状態の境界です。ルールはUpdateへ足し、Drawの順番や見た目は自由に差し替えられます。下の01から順に、その中へ各仕組みを実装していきます。`;
    if(slug === "active-rpg") pair[1] += ` ゲージだけを足すと、同じフレームに満タンになった役を一つの変数で上書きしてしまいます。READYキューなら、満タンになった順を全員ぶん保存してから一人ずつ解決できます。`;
    if(slug === "racing") pair[1] += ` 各ステップは前の状態を引き継ぎ、加速→旋回→ゲート→ライバルの順で同じ車へ重ねます。`;
    if(slug === "metroidvania") pair[1] += ` 下の01〜05もこの順番です。カメラで見る範囲を作ってから、部屋・地図・能力ゲート・戻り道を足します。`;
  }
  // A bridge is useful only when it can explain this lesson's own rule.  A
  // generic Update/Draw card repeated on every page hides the actual lesson.
  if(!pair)continue;
  if(lang==="ja"&&route.includes("/tracks/visual-effects/")){
    pair[1]+=" Go用語メモ：整数型（int）は小数を持たない数、浮動小数点型（float32/float64）は小数を持てる数、typeはデータの箱の種類を決める言葉、funcは処理をひとまとめにする関数の印です。[]TはTを順番に並べるリスト、appendはその末尾へ一つ足す命令です。Updateは数字と状態を進め、Drawはその結果を絵にします。";
  }
  if(lang==="ja"&&route.includes("/tracks/")&&!route.includes("/tracks/visual-effects/")&&!isHub){
    pair[1]+=" PLAYABLEは『このページで実際に操作できる』という印です。";
  }
  if(lang==="en")pair=["Trace one rule first","Find where this one line is used, then follow how its value changes on screen. This is the shared game-state boundary: add rules in Update, then choose an independent presentation in Draw. The numbered steps below add one system at a time.",pair[2] || "func (g *Game) Update() error { return nil }"];
  pair[2] ||= "func (g *Game) Update() error { return nil }";
  const esc=s=>s.replaceAll("&","&amp;").replaceAll("<","&lt;");
  const flowVisual = isHub && trackParts[slug]
    ? `<div class="course-mini-flow" role="img" aria-label="${lang === "ja" ? "コースの学習順" : "Course learning order"}">${trackParts[slug].split("→").map((step, i) => `<span>${String(i + 1).padStart(2, "0")} ${step}</span>`).join('<b aria-hidden="true">→</b>')}</div>`
    : "";
  const block=`<!-- BEGIN BEGINNER BRIDGE -->\n<section class="beginner-bridge"><div><p class="eyebrow">${lang==="ja"?"はじめて読む人へ":"FIRST-TIME READER"}</p><h2>${pair[0]}</h2><p>${pair[1]}</p></div><div class="beginner-code"><p><strong>${lang==="ja"?"このページの1本":"ONE RULE TO TRACE"}</strong><br>${lang==="ja"?"説明と次のコードを対応させてから、ゲームで変化を確かめよう。":"Match the explanation to this line, then test the change in the game."}</p><pre><code>${esc(pair[2])}</code></pre></div>${flowVisual}</section>\n<!-- END BEGINNER BRIDGE -->\n`;

  // A playable lesson must keep its promise: let the reader play before
  // introducing vocabulary or code. Put the orientation immediately after
  // the first play section; hubs without a game still use the hero.
  const playStart = html.search(/<section[^>]*class=["'][^"']*play[^"']*["'][^>]*>/i);
  if (playStart >= 0) {
    const playEnd = html.indexOf("</section>", playStart);
    if (playEnd >= 0) {
      const at = playEnd + "</section>".length;
      html = html.slice(0, at) + "\n" + block + html.slice(at);
      writeFileSync(path, html);
      continue;
    }
  }
  // Hubs and setup guides have no playable section, so orient the reader
  // immediately after their introduction.
  const heroClasses=["data-hero","lesson-hero","track-hero","overview-hero","setup-hero","catalog-hero"];
  let inserted=false;
  for(const className of heroClasses){
    const start=html.search(new RegExp(`<section[^>]*class=["'][^"']*\\b${className}\\b[^"']*["']`));
    if(start<0)continue;
    const end=html.indexOf("</section>",start);
    if(end<0)continue;
    const at=end+"</section>".length;
    html=html.slice(0,at)+"\n"+block+html.slice(at);
    inserted=true;
    break;
  }
  if(!inserted){
    const anchor=html.includes('<section class="path-list"')?'<section class="path-list"':html.includes('<section class="play-panel"')?'<section class="play-panel"':html.includes('<section class="feedback-section"')?'<section class="feedback-section"':"</main>";
    html=html.replace(anchor,block+anchor);
  }
  writeFileSync(path,html);
}
