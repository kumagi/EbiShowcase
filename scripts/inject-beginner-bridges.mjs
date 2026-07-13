#!/usr/bin/env node
import { readdirSync, readFileSync, statSync, writeFileSync } from "node:fs";
import { join, relative } from "node:path";

const root = new URL("../web", import.meta.url).pathname;
const marker = /<!-- BEGIN BEGINNER BRIDGE -->[\s\S]*?<!-- END BEGINNER BRIDGE -->\s*/g;

// Short, concrete bridges requested by classroom review.  These complement the
// deep dive: first define the unfamiliar word, then show one rule to trace.
const notes = {
  "home":"ゲームループは『数字を進める→その数字で絵を描く』を1秒に約60回くり返す仕組み。海老・天次郎（えび・てんじろう）が的を追う最小例から始めます。|func (g *Game) Update() error { g.x++; return nil }",
  "setup":"go.modは『このゲームの名前と使う部品』を書いた目次ファイル。PowerShell/ターミナルは文字で命令する同じ役目の窓です。天次郎の最初の一コマを出すところまで進みます。|func (g *Game) Update() error { return nil }",
  "game-data":"データ駆動とは、村・台詞・敵をコードへ直接書かず、交換できるデータとして読む作り方。asset loaderはそのファイルをゲームへ渡す係、IDは部品同士を結ぶ名前札です。|village := loadMap(\"village.json\")",
  "tap-target":"メソッドはGameという箱に結び付いた関数。Updateは天次郎や的の数字を変え、Drawは現在の数字を絵にし、Layoutは画面の基準サイズを返します。|func (g *game) Update() error { return nil }",
  "timing-meter":"状態は『狙っている・判定を見ている・再挑戦待ち』のような今の場面。switchは状態ごとの処理をif/elseより読みやすく並べます。%は割り算の余りです。|phase := frame % 120",
  "catch-stars":"スライスは星を何個でも入れられる順番付きリスト。天次郎が星を受け止め、3個落とすとゲームオーバーです。|stars = append(stars, newStar())",
  "flappy":"加速度は『速さが毎フレームどれだけ変わるか』。羽ばたき時の-7.4は上向きの初速度で、重力が少しずつ下向きへ戻します。|vy = -7.4; vy += gravity; y += vy",
  "pong":"ベクトルは横の速さvxと縦の速さvyを組にしたもの。符号反転は+を-へ変え、壁へ向かう速さを反対向きにします。|vx = -vx",
  "breakout":"[]は順番付きの箱、[][]は行と列の方眼紙。全部のbrickが消えたらwinがclear状態を立て、次のDrawで結果画面を出します。|if remaining == 0 { state = Win }",
  "snake":"スライスは長さを変えられるリスト。ヘビの体は通ったマス座標の順番で、appendで頭を足し、最後を外すと一歩進みます。|body = append([]Cell{head}, body[:len(body)-1]...)",
  "space-shooter":"自機・弾・敵を別リストにすると、それぞれ同じUpdateを繰り返せます。全弾×全敵を調べる総当たりは、教材の小さな数なら十分です。|for _, b := range bullets { for _, e := range enemies { check(b,e) } }",
  "sokoban":"論理座標は『何列・何行』、描画座標は画面のピクセル。天次郎が1マス動く間はmoving=trueにして入力を待たせ、Drawだけ途中位置を描きます。|if moving { progress++; return nil }",
  "platformer":"めり込みは物体が壁の中へ重なること。Xを動かして直し、次にYを動かして直すと角でも滑れます。0.09はカメラが差の9%ずつ追う調整値です。|cameraX += (playerX-cameraX) * .09",
  "dungeon":"patrol（巡回）は決まった道の見回り、chase（追跡）は天次郎を追う状態。165pxなどは気付く距離を遊びながら調整した値です。|if distance < 165 { state = Chase }",
  "bullet-hell":"cos/sinは角度から円周上の横・縦位置を得る道具。ラジアンは角度単位で2πが一周、frame%周期は一定フレームごとに同じ模様を繰り返します。|angle := 2*math.Pi*float64(i)/count",
  "moving-platforms":"床自身のxを毎フレーム動かし、天次郎が上に立っている間は床の移動量dxを天次郎にも足します。台詞は画面下の案内欄へ描きます。|platform.x += platform.vx; player.x += platform.vx",
  "patrol-enemies":"patrol AIは二つの端点を往復する見回り役。端へ着いたら向きを反転し、天次郎が近ければ追跡状態へ変えます。|if enemy.x <= left || enemy.x >= right { enemy.vx = -enemy.vx }",
  "powerup-adventure":"取得矩形が重なるとspeed・jump・HPのどれを変えるかをitem.kindで選びます。何が強くなったかを数値と色で同時に表示します。|if playerRect.Overlaps(itemRect) { apply(item.kind) }",
  "scrolling-stage":"worldXは世界での位置、screenX=worldX-cameraXが画面位置。カメラは天次郎が中央を越えた分だけ追います。|screenX := worldX - cameraX",
  "tiny-platformer":"最小の三部品は天次郎、床、ジャンプ。入力でvx/vyを変え、床判定で止め、Drawで二つを描けば遊びになります。|vy += gravity; if onFloor { vy = 0 }",
  "arena-dodge":"この段階の主役は回避です。敵だけが天次郎へ近づき、自動攻撃は次のSTEPで追加します。敵数と速度をstageデータへ分けます。|enemyCount, enemySpeed := 6, 1.2",
  "auto-turret":"turretは一番近い敵を探し、cooldownが0の時だけ弾を作る自動砲台。atan2は敵へのdx,dyを向きへ変える道具です。|if cooldown == 0 { fire(nearestEnemy()) }",
  "experience-draft":"敵を倒すとxpを足し、xp>=nextXPなら時間を止めて3択状態へ移ります。選択後に必要xpを増やします。|if xp >= nextXP { state = ChooseUpgrade }",
  "survival-run":"目標時間まで生き残るゲーム。接触でHPが減り、敵を倒すと経験値、強化選択で攻撃範囲や速さが変わります。|if timer >= goal { state = Clear }",
  "swarm":"全敵を一度ずつUpdateするので敵数に比例して仕事が増えます。画面外描画を省き、上限を決めて60fpsを守ります。|for i := range enemies { enemies[i].Update() }",
  "weapon-evolution":"経験値ゲージが満ちるたびlevelを増やし、level 2で3way、3で連射など武器データを交換します。|weapon = evolutions[level]",
  "command-battle":"A/B/Cまたは画面ボタンをAttack/Guard/Skillへ翻訳し、選んだcommandでダメージ式を分けます。|switch command { case Attack: damage=atk-def }",
  "dialogue-flags":"flagは選択を覚える名前付きのtrue/false。Gameの中に定義し、後の会話でif flagを読んで台詞を分けます。|if choice == Help { g.helped = true }",
  "ebi-quest":"worldStateは今の場面、flagsは過去の選択。村→会話→戦闘の順に状態を進め、完成形でも一度に一部品だけ更新します。|switch state { case Village: updateVillage() }",
  "inventory-shop":"itemsは品物のリスト、goldは全取引で共有する一つの数。価格以上のgoldがある時だけ減らしてitemを足します。|if gold >= item.price { gold -= item.price; add(item) }",
  "stats-status":"structはHP/ATK/DEFを一組にする箱。敵も味方も同じ計算へ渡せるので、値の取り違えを減らせます。|type Stats struct { HP, ATK, DEF int }",
  "village-walk":"map[y][x]は一マスの種類、tileSizeは一マスのピクセル数。描画はx*tileSize,y*tileSizeへ変換します。|screenX := tileX * tileSize",
  "world-encounters":"天次郎の矩形とevent領域が重なったら戦闘状態へ移ります。確率遭遇なら0〜99の乱数と割合を比べます。|if player.Overlaps(eventArea) { state = Battle }",
  "tap-counter":"countはGameの中で宣言し、押した瞬間だけ+1。Ebitengineでは押しっぱなしでなくJustPressedを使います。|if inpututil.IsMouseButtonJustPressed(0) { g.count++ }",
  "first-shop":"買うたび価格を base+owned*step で上げると次の目標が少し遠くなります。商品リストなら同じ購入式を使い回せます。|price := basePrice + owned*10",
  "growth-curves":"成長曲線は回数と必要価格の関係。price=base*1.15^ownedなら、買うたび前回の約115%になります。|price := base * math.Pow(1.15, float64(owned))",
  "idle-factory":"工場ごとのperSecondを合計し、60フレームで割った量を毎Update足します。10msの待機処理は使いません。|stock += perSecond / 60",
  "offline-bakery":"終了時刻を保存し、再開時刻との差秒×毎秒生産量を一度だけ足します。長すぎる放置時間には上限を付けます。|gain := time.Since(savedAt).Seconds() * perSecond",
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
  "ebi-merge":"円の組を番号順に一度だけ調べ、合体した二つへremoved印を付けます。新しい円は次のUpdateから調べ、同じ瞬間の二重合体を防ぎます。|if !a.removed && !b.removed { tryMerge(a,b) }",
  "ebi-blocks":"形[][]、回転候補[][][]、現在位置を別に持ちます。[]が一段増えるごとに『形の一覧』が一層増えると読めます。|rotatedShape := rotations[piece][rotation]",
  "ally-effects":"前STEPのParticle型とspawn関数を同じpackageから呼び、仲間の種類だけ色・速さ・個数の引数で変えます。|spawnEffect(hitX, hitY, ally.effect)",
  "circle-battle":"全敵を一度ずつ調べ、中心距離²が半径合計²未満なら命中。敵数が小さい教材ではこのO(n)走査が最も読みやすい方法です。|for i := range enemies { checkCircle(projectile, enemies[i]) }",
  "ebi-strike":"完成形はaim→launch→move→collide→fxの五関数。どのSTEPの部品かが呼び出し順で分かります。|aim(); launch(); move(); collide(); updateFX()",
  "falling-pair":"入力中は列だけ選び、Updateの落下tickで二個のrowを同時に+1。Drawは親と相方の二座標を同じ進み具合で描きます。|pair.y += fallSpeed",
  "ebi-chain":"完成形はinput→fall→lock→findGroup→clear→gravityをphaseで一つずつ進めます。最初はこの一列を復習地図として読みます。|updateInput(); updateFall(); resolveBoard()",
  "vfx-stamp":"GeoM は絵の置き場所を覚える道具です。Translate(x,y) の x は右、y は下へ動かす量。& は設定箱そのものを渡す印です。|op.GeoM.Translate(x, y)",
  "vfx-transform":"w/h は width（横幅）/height（高さ）。-w/2,-h/2 で絵の中心を座標(0,0)へ合わせてから回します。ラジアンは角度の単位で、πが180度です。|op.GeoM.Translate(-w/2, -h/2)",
  "vfx-tint":"tint は元絵へ色を重ねること。ColorScale は赤・緑・青・透明度の順で、1はそのまま、0.4は40%残す意味です。|op.ColorScale.Scale(1, .4, .4, 1)",
  "vfx-additive":"加算合成は、すでにある光へ次の光の明るさを足す描き方。赤い光と緑の光は重なると黄色へ近づきます。|op.Blend = ebiten.BlendLighter",
  "vfx-walk":"SubImage（サブイメージ）は、大きな画像から四角い一部分を切り出す命令。タイマーは何コマ進んだかを数える整数（小数でない数）です。type は、その数や絵をまとめる新しい箱の種類を決める書き方。col は横、row は縦の何番目かで、col×コマ幅が左端になります。Updateがコマ番号を変え、Drawが切り出した絵を表示します。加算合成は光を足す描き方、パーティクルは短時間だけ動く小さな粒です。|frame := sheet.SubImage(frameRect).(*ebiten.Image)",
  "vfx-particles":"particle は一粒分の位置・速さ・残り時間をまとめた箱。全粒のlifeを減らし、0より大きい粒だけ次のリストへ残します。|if p.life > 0 { next = append(next, p) }",
  "vfx-spells":"エフェクトの合成は、別々に作った光・魔法本体・火花を順番に重ねて一つに見せること。ScaleAlphaは透明度を何倍にするか、ColorScaleは赤・緑・青・透明度を何倍にするかの設定です。整数型は小数を持たない数の種類、typeは新しいデータの種類を定義するGoの言葉。Updateが粒の数字を進め、Drawが背景の光→本体→粒の順に描きます。パーティクルは生まれて動いて消える小さな粒です。|drawGlow(); drawSpell(); drawParticles()",
  "vfx-magic-fire":"押しているフレーム数がチャージ量。しきい値は段階が変わる境目で、60フレームなら約1秒です。|charge++ // 60 frames is about 1 second",
  "vfx-magic-ice":"phaseを形成と粉砕に分け、前半は枝を伸ばし、後半は同じ枝を破片へ分けます。二段階だから育って砕けたように読めます。|if phase == Form { grow() } else { shatter() }",
  "vfx-magic-thunder":"再帰は、線を分ける処理が自分自身をもう一度呼ぶ書き方。角度は見た目の調整値なので、まず±30度から比べます。|branch(a, b, depth-1)",
  "vfx-magic-light":"16本は円を22.5度ずつ分けた見栄えの例。glow画像を薄く重ねると輪郭がぼけた光に見えます。|angle := 2 * math.Pi * float64(i) / 16",
  "vfx-magic-dark":"らせんは角度を増やしながら中心までの距離rを減らす点の列。吸い込み強さ2.0は一度に近づく量です。|r -= 2; angle += .12",
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
  "clear-and-fall":"各列を下から上へ読み、空でないピースだけを下のwrite行へ移します。要素を削除せず上へ空マスを集める方法です。|for y := h-1; y >= 0; y-- { compact(y) }",
  "find-matches":"横はy固定でxを、縦はx固定でyを進めます。隣を見る前にx+1<幅、y+1<高さを確かめます。|if x+1 < gridW && grid[y][x] == grid[y][x+1]",
  "grid-swap":"Goは a,b=b,a で二つを同時交換できます。交換後に揃いがなければ同じ式でもとへ戻します。|a, b = b, a",
  "special-pieces":"4個なら行/列消去、5個なら同色消去というkindをピースに保存し、通常消去とは別関数で範囲を決めます。|piece.kind = RowClear",
  "box-viewer":"赤いattack boxと青いhurt boxを画像へ重ねる確認画面。画像ではなく二つの四角が重なったかで命中を決めます。|hit := attack.Overlaps(hurt)",
  "combo-cancel":"攻撃の特定フレームだけ次の入力を予約します。15フレームは60fpsで約0.25秒です。|if frame >= 8 && frame <= 15 { acceptNext() }",
  "command-fighter":"押した種類とframeを履歴へ追加し、20フレームより古い入力を捨ててから並びを比べます。|inputs = append(inputs, Input{key, frame})",
  "frame-attack":"60fpsで4フレームは約0.067秒。攻撃矩形はactiveの間だけ出し、その前後は当たりません。|active := frame >= 6 && frame < 10",
  "guard-throw":"投げ入力→ガード→通常移動の優先順で一つだけ状態を選び、同時に二状態を立てません。|switch { case throw: state=Throw; case guard: state=Guard }",
  "hit-reaction":"攻撃矩形と食らい矩形が重なった一瞬にHPを減らし、同じ出来事からのけぞり演出を始めます。|if attack.Overlaps(hurt) { hp--; state=Hurt }",
  "circle-collision":"中心間の距離が半径A+半径Bより小さいと円は重なります。平方根なしで距離²同士を比べても同じです。|dx*dx+dy*dy < (ra+rb)*(ra+rb)",
  "falling-circles":"g=.5は毎フレーム下向き速度が.5増える教材値。円の下端y+rが床を越えたら床-rへ戻します。|vy += .5; y += vy",
  "merge-rule":"同じtierの円が触れたら二つを消し、中心へtier+1を一つ作ります。|if a.tier == b.tier { merge(a, b) }",
  "stable-simulation":"一度で強く押し戻さず、めり込み修正を4回に分けます。4は速さと安定の折衷値で、変えて比較できます。|for pass:=0; pass<4; pass++ { solveContacts() }",
  "stacking-bounce":"反発係数.8なら衝突前の速さの80%で反対へ戻り、20%を失います。|vy = -vy * .8",
  "bag-hold-ghost":"7-bagは7種類を一つずつ袋へ入れ、順番だけ混ぜる方法。同じ形ばかり出る事故を防ぎます。|rng.Shuffle(len(bag), swap)",
  "tetromino-shapes":"外側の[]が行、内側の[]が列。1のマスだけ描けば同じ仕組みで7形を表せます。|if shape[row][col] == 1 { drawBlock() }",
  "rotation-kicks":"回転後に重なるなら0、左1、右1、左2、右2の位置を試し、全部だめなら回転を戻します。|for _, dx := range []int{0,-1,1,-2,2} { try(dx) }",
  "lock-lines":"満杯行より上を一段ずつ下へコピーし、最上段を空にすると、難しいslice削除を使わずに済みます。|copy(grid[1:y+1], grid[0:y])",
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
  "buffered-turn":"毎フレーム入力を見てwantedDirへ予約し、次の交差点で通れる時に採用します。|wantedDir = pressedDir",
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
  [/Update|Draw/,"Updateは次の一コマの数字を決め、Drawはその数字を絵にします。これがEbitengineのゲームループです。"],
];

function walk(dir){return readdirSync(dir).flatMap(n=>{const p=join(dir,n);return statSync(p).isDirectory()?walk(p):[p]})}
for(const path of walk(root).filter(p=>p.endsWith(".html"))){
  let html=readFileSync(path,"utf8").replace(marker,"");
  const route="/"+relative(root,path).replace(/index\.html$/,"").replaceAll("\\","/");
  if(!html.includes("<main")||!/(\/games\/|\/tracks\/|\/guides\/|^\/(ja|en)\/$)/.test(route))continue;
  const bits=route.split("/").filter(Boolean), lang=bits[0], slug=bits.length===1?"home":bits.at(-1);
  const isHub=bits[1]==="tracks"&&bits.length===3;
  let pair=notes[slug]?["ここを先に見よう",...notes[slug].split("|")]:undefined;
  if(isHub&&trackParts[slug]) pair=[`このコースの組み立て図`,`${trackParts[slug]}の順に、海老・天次郎（えび・てんじろう）と一部品ずつ作ります。`,`func (g *Game) Update() error { return nil }`];
  if(!pair){const hit=generic.find(([r])=>r.test(html)); if(hit)pair=["コードを読む前の用語メモ",hit[1],"func (g *Game) Update() error { return nil }"]}
  if(!pair)continue;
  if(lang==="en")pair=["Read the game as small parts","Update changes Ebi Tenjiroh's game state; Draw shows that state as the next frame. In func (g *Game), (g *Game) is the receiver—the name tag for this Game.",pair[2] || "func (g *Game) Update() error { return nil }"];
  pair[2] ||= "func (g *Game) Update() error { return nil }";
  const esc=s=>s.replaceAll("&","&amp;").replaceAll("<","&lt;");
  const block=`<!-- BEGIN BEGINNER BRIDGE -->\n<section class="beginner-bridge"><div><p class="eyebrow">${lang==="ja"?"はじめて読む人へ":"FIRST-TIME READER"}</p><h2>${pair[0]}</h2><p>${pair[1]}</p></div><div class="beginner-code"><p><strong>${lang==="ja"?"海老・天次郎の見方":"Ebi Tenjiroh's view"}</strong><br>${lang==="ja"?"Updateが天次郎の数字を変え、Drawが次の一コマを描きます。":"Update changes Tenjiroh's numbers; Draw paints the next frame."}</p><pre><code>${esc(pair[2])}</code></pre><p class="receiver-note">${lang==="ja"?"func (g *Game)の(g *Game)は、この仕事がGameの値を読むための名札（レシーバー）です。":"(g *Game) is the receiver: a name tag that lets the function use this Game's state."}</p></div></section>\n<!-- END BEGINNER BRIDGE -->\n`;
  const anchor=html.includes('<section class="path-list"')?'<section class="path-list"':html.includes('<section class="play-panel"')?'<section class="play-panel"':"</main>";
  writeFileSync(path,html.replace(anchor,block+anchor));
}
