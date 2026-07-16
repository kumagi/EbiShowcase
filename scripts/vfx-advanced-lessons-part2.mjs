/** Advanced VFX lessons A07–A12 (continuation). */
export const advancedLessonsPart2 = [
  {
    slug: "vfx-fx-snake",
    step: "A07",
    tier: "advanced",
    core: "snake",
    stars: "★★★★★",
    labKind: "fx-snake",
    concept: { ja: "体セルとキラキラ", en: "Body cells vs sparkles" },
    hubDesc: {
      ja: "LEVEL 07 Snake。体は []cell、食べた瞬間だけ fx。",
      en: "LEVEL 07 Snake—body []cell, FX only on eat.",
    },
    ja: {
      navConcept: "履歴と演出",
      title: "スネークに捕食演出を足す",
      lead: "体はマス目の履歴。食べた瞬間の花火だけ fx。体の各セルに粒を持たせない。",
      deepEyebrow: "DEEP DIVE / GRID BODY",
      deepH: "体は履歴、<br>花火は一瞬",
      deepLead: "Snake の本質は位置の列。演出をセルに埋め込むと、移動ロジックが読めなくなります。",
      concepts: [
        { h: "体", p: "グリッド上のセル列。", code: "body []cell" },
        { h: "捕食", p: "食った座標で Burst。", code: "fx.Burst(food)" },
        { h: "分離", p: "セルに particle を持たない。", code: "no cell.fx" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "セルはデータ", p: "見た目のキラキラは体データに混ぜない。" },
      play: { title: "餌を食べて伸びよう", p: "方向を変えて餌へ。捕食花火は fx。body と fx.parts を見比べよう。" },
      codeHead: { eyebrow: "IN THE DEMO", h: "eat の一箇所だけ", p: "移動ループの中で毎セル Burst しない。" },
      whys: [
        { eyebrow: "WHY?", h: "グリッドが本体", p: "学習ポイントは履歴。演出は添え物。" },
        { eyebrow: "WHY ONE BURST?", h: "イベント", p: "食べた瞬間だけ祝う。" },
        { eyebrow: "TRY NEXT", h: "シューターへ", p: "弾・敵・演出の三分割。" },
      ],
    },
    en: {
      navConcept: "History vs juice",
      title: "Snake Eat FX",
      lead: "Body is grid history. Only the eat moment gets fx. Don't hang particles on every cell.",
      deepEyebrow: "DEEP DIVE / GRID BODY",
      deepH: "Body is history,<br>fireworks are a moment",
      deepLead: "Snake is a list of cells. Embedding FX per cell muddies the move logic.",
      concepts: [
        { h: "Body", p: "Cells on a grid.", code: "body []cell" },
        { h: "Eat", p: "Burst at food.", code: "fx.Burst(food)" },
        { h: "Split", p: "No particle on cell.", code: "no cell.fx" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "Cells are data", p: "Don't mix sparkles into body data." },
      play: { title: "Eat and grow", p: "Steer to food. Eat fireworks are fx—compare body vs fx.parts." },
      codeHead: { eyebrow: "IN THE DEMO", h: "One place: eat", p: "Don't Burst every cell while moving." },
      whys: [
        { eyebrow: "WHY?", h: "Grid is the lesson", p: "History matters; juice is garnish." },
        { eyebrow: "WHY ONE BURST?", h: "Event", p: "Celebrate the eat instant only." },
        { eyebrow: "TRY NEXT", h: "Shooter", p: "Bullets, enemies, FX—three piles." },
      ],
    },
    code: `// on eat:
body = grow(body)
g.fx.Burst(food.x, food.y, 24, 3, tint, true)`,
  },
  {
    slug: "vfx-fx-shooter",
    step: "A08",
    tier: "advanced",
    core: "space-shooter",
    stars: "★★★★★",
    labKind: "fx-split",
    concept: { ja: "弾・敵・FXの三分割", en: "Bullets / enemies / FX" },
    hubDesc: {
      ja: "LEVEL 08 シューター。弾・敵・マズルフラッシュを分ける。",
      en: "LEVEL 08 shooter—bullets, enemies, muzzle FX apart.",
    },
    ja: {
      navConcept: "三種のリスト",
      title: "弾と敵と爆発を分ける",
      lead: "自弾、敵、爆発粒。役割ごとにスライス。撃った瞬間のマズルは fx、敵の HP は play。",
      deepEyebrow: "DEEP DIVE / THREE LISTS",
      deepH: "当たるもの／倒すもの／見えるだけ",
      deepLead: "弾と敵は当たり判定あり。爆発に当たり判定は不要。三分割が見えてきたら応用編のゴールに近いです。",
      concepts: [
        { h: "弾", p: "動いて当たる。", code: "bullets []" },
        { h: "敵", p: "HP と出現。", code: "enemies []" },
        { h: "FX", p: "マズルと爆発。", code: "fx.Burst" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "三つの箱", p: "弾・敵・演出。どれが欠けても別の壊れ方をする。" },
      play: { title: "撃って爆発させよう", p: "移動とショット。敵撃破の爆発は fx。リスト数を LIVE GO で確認。" },
      codeHead: { eyebrow: "IN THE DEMO", h: "kill したら削除＋爆発", p: "敵を fx 粒子に格下げしない。" },
      whys: [
        { eyebrow: "WHY THREE?", h: "更新が違う", p: "弾は直線、敵はAI、FXは寿命。" },
        { eyebrow: "WHY NOT CONVERT?", h: "型の谎言", p: "敵を粒に変えると当たり判定のバグが残る。" },
        { eyebrow: "TRY NEXT", h: "倉庫番へ", p: "タイル地図と押し演出。" },
      ],
    },
    en: {
      navConcept: "Three lists",
      title: "Bullets, Enemies, Explosions",
      lead: "Player shots, enemies, and burst FX. Muzzle flash is fx; HP is play.",
      deepEyebrow: "DEEP DIVE / THREE LISTS",
      deepH: "What hits / what dies / what only shows",
      deepLead: "Bullets and enemies collide. Explosions shouldn't. Three piles means you're near the advanced goal.",
      concepts: [
        { h: "Bullets", p: "Move and hit.", code: "bullets []" },
        { h: "Enemies", p: "HP and spawns.", code: "enemies []" },
        { h: "FX", p: "Muzzle and blasts.", code: "fx.Burst" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "Three boxes", p: "Bullets, enemies, FX—each breaks differently if missing." },
      play: { title: "Shoot and explode", p: "Move and fire. Death blasts are fx—check counts in LIVE GO." },
      codeHead: { eyebrow: "IN THE DEMO", h: "On kill: delete + burst", p: "Don't demote enemies into particles." },
      whys: [
        { eyebrow: "WHY THREE?", h: "Different updates", p: "Shots, AI, lifetimes." },
        { eyebrow: "WHY NOT CONVERT?", h: "Type lies", p: "Turning foes into sparks leaves hitbox bugs." },
        { eyebrow: "TRY NEXT", h: "Sokoban", p: "Tile map vs push dust." },
      ],
    },
    code: `type Game struct {
  bullets []bullet
  enemies []enemy
  fx      vfxfx.System
}`,
  },
  {
    slug: "vfx-fx-sokoban",
    step: "A09",
    tier: "advanced",
    core: "sokoban",
    stars: "★★★★★",
    labKind: "fx-split",
    concept: { ja: "タイル地図と押し演出", en: "Tilemap vs push FX" },
    hubDesc: {
      ja: "LEVEL 09 倉庫番。地図は数値、押しは fx の塵。",
      en: "LEVEL 09 sokoban—numeric map, push dust in fx.",
    },
    ja: {
      navConcept: "データと汁",
      title: "箱押しに塵とゴール輝きを足す",
      lead: "ステージは二次元の数値。押した瞬間の塵、ゴールの輝きは fx。タイルに演出フラグを増やしすぎない。",
      deepEyebrow: "DEEP DIVE / TILEMAP",
      deepH: "地図は真実、<br>塵は気持ち",
      deepLead: "ルール判定はタイル値だけ見てよい。見た目のキラキラで判定するとバグの温床になります。",
      concepts: [
        { h: "地図", p: "壁・床・箱・ゴール。", code: "tiles [][]int" },
        { h: "押し", p: "成功したら塵。", code: "fx.Burst on push" },
        { h: "ゴール", p: "箱が乗ったら輝き。", code: "fx.Ring" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "判定はタイル", p: "キラキラの有無でクリア判定しない。" },
      play: { title: "箱をゴールへ", p: "箱を押す。塵と輝きは fx。地図データは play。" },
      codeHead: { eyebrow: "IN THE DEMO", h: "clear はタイル比較", p: "FX は結果の演出だけ。" },
      whys: [
        { eyebrow: "WHY?", h: "セーブとテスト", p: "地図だけ保存すれば再現できる。" },
        { eyebrow: "WHY FX OPTIONAL?", h: "ルールが先", p: "演出なしでもパズルは解ける。" },
        { eyebrow: "TRY NEXT", h: "プラットフォーマーへ", p: "接地とジャンプ塵。" },
      ],
    },
    en: {
      navConcept: "Data vs juice",
      title: "Sokoban Dust & Goal Glow",
      lead: "The stage is numbers. Push dust and goal sparkles are fx—don't overload tiles with FX flags.",
      deepEyebrow: "DEEP DIVE / TILEMAP",
      deepH: "Map is truth,<br>dust is feeling",
      deepLead: "Rules should read tile values only. Judging by sparkles invites bugs.",
      concepts: [
        { h: "Map", p: "Wall, floor, box, goal.", code: "tiles [][]int" },
        { h: "Push", p: "Dust when push succeeds.", code: "fx.Burst on push" },
        { h: "Goal", p: "Glow when box sits.", code: "fx.Ring" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "Judge tiles", p: "Don't clear based on sparkles." },
      play: { title: "Push boxes to goals", p: "Dust and glow are fx; the map is play." },
      codeHead: { eyebrow: "IN THE DEMO", h: "Clear = tile compare", p: "FX only celebrates the result." },
      whys: [
        { eyebrow: "WHY?", h: "Save & test", p: "Persist the map and you can replay." },
        { eyebrow: "WHY FX OPTIONAL?", h: "Rules first", p: "Puzzle works with juice off." },
        { eyebrow: "TRY NEXT", h: "Platformer", p: "Landing and jump dust." },
      ],
    },
    code: `if pushOK {
  tiles = applyPush(tiles) // play truth
  g.fx.Burst(boxX, boxY, 16, 2, dust, true)
}`,
  },
  {
    slug: "vfx-fx-platform",
    step: "A10",
    tier: "advanced",
    core: "platformer",
    stars: "★★★★★",
    labKind: "fx-split",
    concept: { ja: "接地イベントと塵", en: "Landing events & dust" },
    hubDesc: {
      ja: "LEVEL 10 横アクション。ジャンプ／着地の塵は fx。",
      en: "LEVEL 10 platformer—jump/land dust in fx.",
    },
    ja: {
      navConcept: "接地フラグ→演出",
      title: "ジャンプと着地に塵を足す",
      lead: "wasOnGround から onGround へ切り替わった瞬間だけ着地演出。物理の grounded は play、塵は fx。",
      deepEyebrow: "DEEP DIVE / EDGE TRIGGER",
      deepH: "tickごとではなく、<br>切り替わった瞬間",
      deepLead: "着地中ずっと Burst すると煙モクモク。エッジ（立ち上がり）で一回だけ呼ぶのもイベント駆動です。",
      concepts: [
        { h: "物理", p: "速度・足場・接地。", code: "onGround bool" },
        { h: "エッジ", p: "false→true で着地。", code: "landed := !was && on" },
        { h: "FX", p: "その瞬間だけ塵。", code: "fx.Burst" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "エッジで一回", p: "押しっぱなし＝連発、ではない。" },
      play: { title: "跳んで着地しよう", p: "ジャンプと着地の塵は fx。接地判定は play。" },
      codeHead: { eyebrow: "IN THE DEMO", h: "前のtickと比較", p: "状態の差分がイベントになる。" },
      whys: [
        { eyebrow: "WHY EDGE?", h: "連打防止", p: "接地中はtickごとにイベントにしない。" },
        { eyebrow: "WHY?", h: "物理が主", p: "プラットフォーマーの核は移動。" },
        { eyebrow: "TRY NEXT", h: "ダンジョンへ", p: "攻撃ヒットの閃光。" },
      ],
    },
    en: {
      navConcept: "Grounded → FX",
      title: "Jump & Land Dust",
      lead: "Rising edge wasOnGround to onGround triggers land FX. Grounded flag is play; dust is fx.",
      deepEyebrow: "DEEP DIVE / EDGE TRIGGER",
      deepH: "Not every tick—<br>the moment it changes",
      deepLead: "Bursting while grounded forever makes a smoke factory. Fire once on the edge.",
      concepts: [
        { h: "Physics", p: "Vel, platforms, grounded.", code: "onGround bool" },
        { h: "Edge", p: "false→true = land.", code: "landed := !was && on" },
        { h: "FX", p: "Dust only then.", code: "fx.Burst" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "Once per edge", p: "Held is not repeating events." },
      play: { title: "Jump and land", p: "Jump/land dust is fx; grounded tests are play." },
      codeHead: { eyebrow: "IN THE DEMO", h: "Compare last frame", p: "State deltas become events." },
      whys: [
        { eyebrow: "WHY EDGE?", h: "No spam", p: "Don't event every grounded frame." },
        { eyebrow: "WHY?", h: "Physics first", p: "Movement is the genre core." },
        { eyebrow: "TRY NEXT", h: "Dungeon", p: "Hit flashes on attack." },
      ],
    },
    code: `landed := !g.wasGround && g.onGround
g.wasGround = g.onGround
if landed { g.fx.Burst(g.x, g.y+h, 20, 2.5, dust, true) }`,
  },
  {
    slug: "vfx-fx-dungeon",
    step: "A11",
    tier: "advanced",
    core: "dungeon",
    stars: "★★★★★",
    labKind: "fx-split",
    concept: { ja: "ヒット反応はFX", en: "Hit reactions are FX" },
    hubDesc: {
      ja: "LEVEL 11 ダンジョン。ダメージ計算とヒット閃光を分離。",
      en: "LEVEL 11 dungeon—damage math vs hit flash.",
    },
    ja: {
      navConcept: "HPと閃光",
      title: "攻撃ヒットに閃光を足す",
      lead: "HP を減らすのは play。閃光と破片は fx。ノックバックを入れるなら物理か演出か、意図を決めてから場を選ぶ。",
      deepEyebrow: "DEEP DIVE / COMBAT FEEDBACK",
      deepH: "痛い数字と、<br>痛い見た目",
      deepLead: "数値の真実（HP）とフィードバック（フラッシュ）を分けると、無敵時間やコンボ演出を足しやすくなります。",
      concepts: [
        { h: "ダメージ", p: "HP を減らす。", code: "hp -= dmg" },
        { h: "反応", p: "フラッシュと粒。", code: "fx.FlashScreen" },
        { h: "分離", p: "見た目だけで倒した判定にしない。", code: "bar from hp" },
      ],
      lab: { eyebrow: "TRY IT / PLAY vs FX", title: "命中＝数字＋花火、どちらも別リスト", p: "「命中！」で play の点数と fx の粒が同時に増える。「fx 1F」では粒だけ減り、点数は残る。" },
      play: { title: "スライムを叩こう", p: "移動と攻撃。ヒット演出は fx。HP は play。" },
      codeHead: { eyebrow: "IN THE DEMO", h: "applyDamage の末尾で FX", p: "先に数字、あとで見た目。" },
      whys: [
        { eyebrow: "WHY?", h: "デバッグ", p: "HP ログと演出を切り分けられる。" },
        { eyebrow: "WHY FLASH?", h: "フィードバック", p: "当たった感がないと下手に感じる。" },
        { eyebrow: "TRY NEXT", h: "弾幕へ", p: "総仕上げ。弾クリアと大爆発。" },
      ],
    },
    en: {
      navConcept: "HP vs flash",
      title: "Dungeon Hit Flashes",
      lead: "Reducing HP is play. Flashes and shards are fx. Choose knockback as physics or juice on purpose.",
      deepEyebrow: "DEEP DIVE / COMBAT FEEDBACK",
      deepH: "Painful numbers vs<br>painful looks",
      deepLead: "Split HP truth from flash feedback and i-frames/combos get easier to add.",
      concepts: [
        { h: "Damage", p: "Reduce HP.", code: "hp -= dmg" },
        { h: "Reaction", p: "Flash and sparks.", code: "fx.FlashScreen" },
        { h: "Split", p: "Don't invent HP from blinks.", code: "bar from hp" },
      ],
      lab: { eyebrow: "TRY IT / PLAY vs FX", title: "A hit updates numbers and sparks", p: "Hit! grows score + sparks. Tick FX ages sparks only — score stays." },
      play: { title: "Smack the slime", p: "Move and attack. Hit juice is fx; HP is play." },
      codeHead: { eyebrow: "IN THE DEMO", h: "FX at end of applyDamage", p: "Numbers first, looks second." },
      whys: [
        { eyebrow: "WHY?", h: "Debug", p: "Log HP without the light show." },
        { eyebrow: "WHY FLASH?", h: "Feedback", p: "No hit feel = feels unfair." },
        { eyebrow: "TRY NEXT", h: "Bullet hell", p: "Finale: clear bombs vs bullets[]." },
      ],
    },
    code: `func applyDamage(e *enemy, dmg int) {
  e.hp -= dmg
  g.fx.FlashScreen(0.5, 255, 80, 80)
  g.fx.Burst(e.x, e.y, 28, 4, tint, true)
}`,
  },
  {
    slug: "vfx-fx-bullethell",
    step: "A12",
    tier: "advanced",
    core: "bullet-hell",
    stars: "★★★★★",
    labKind: "fx-split",
    concept: { ja: "弾幕クリアとFXの総仕上げ", en: "Bomb clear vs bullets[]" },
    hubDesc: {
      ja: "LEVEL 12 弾幕。弾配列を消し、見た目の爆発は fx。",
      en: "LEVEL 12 bullet hell—delete bullets[], explode in fx.",
    },
    ja: {
      navConcept: "総仕上げ",
      title: "ボムで弾を消し、爆発は演出へ",
      lead: "応用編のゴール。ボムは近くの bullets を削除（play）し、同時に巨大な Ring と Flash（fx）。弾を演出用に再利用しない——消してから花火。",
      deepEyebrow: "DEEP DIVE / CAPSTONE",
      deepH: "危険な配列と、<br>安全な花火",
      deepLead: "bullets は当たるとダメージ。fx.Parts は当たっても痛くない。この区別を自分の言葉で説明できたら、エフェクト設計の第一関門クリアです。",
      concepts: [
        { h: "弾", p: "危険。当たり判定あり。", code: "bullets []" },
        { h: "ボム", p: "弾を削除するルール。", code: "filter bullets" },
        { h: "花火", p: "安全。見た目だけ。", code: "fx.Ring+Burst" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "危険／安全", p: "左＝弾（危険）、右＝粒（安全）。" },
      play: {
        title: "避けてボムせよ",
        p: "弾幕を避け、ボムで一掃。削除は play、爆発は fx。LIVE GO で bullets と fx.parts を見比べよう。",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "delete してから祝う", p: "Breakout の破壊二段と同じ型の総仕上げ。" },
      whys: [
        { eyebrow: "WHY CAPSTONE?", h: "全部つながる", p: "イベント、状態、配列分割、削除＋FX。" },
        { eyebrow: "WHY SAFE FX?", h: "理不尽を避ける", p: "クリア演出に当たって死ぬと泣ける。" },
        { eyebrow: "TRY NEXT", h: "自分のゲームへ", p: "ジャンルトラックで、同じ分け方を持ち込もう。" },
      ],
    },
    en: {
      navConcept: "Capstone",
      title: "Bomb Clears Bullets—Explosions Are FX",
      lead: "Advanced finale. A bomb deletes nearby bullets (play) and spawns a huge Ring/Flash (fx). Don't recycle live bullets as fireworks—delete, then celebrate.",
      deepEyebrow: "DEEP DIVE / CAPSTONE",
      deepH: "Dangerous arrays vs<br>safe fireworks",
      deepLead: "bullets hurt on contact. fx.Parts never should. If you can explain that split in your own words, you cleared the first gate of effect design.",
      concepts: [
        { h: "Bullets", p: "Dangerous colliders.", code: "bullets []" },
        { h: "Bomb", p: "Rule that deletes them.", code: "filter bullets" },
        { h: "Fireworks", p: "Safe, looks only.", code: "fx.Ring+Burst" },
      ],
      lab: { eyebrow: "TRY IT / SPLIT", title: "Danger / safe", p: "Left = bullets (hurt), right = sparks (safe)." },
      play: {
        title: "Dodge and bomb",
        p: "Survive the spiral, bomb to clear. Deletion is play; explosion is fx—compare bullets vs fx.parts.",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "Delete, then celebrate", p: "Same two-step as Breakout destroy." },
      whys: [
        { eyebrow: "WHY CAPSTONE?", h: "It all connects", p: "Events, state, slice splits, delete+FX." },
        { eyebrow: "WHY SAFE FX?", h: "No cheap deaths", p: "Dying to clear confetti feels awful." },
        { eyebrow: "TRY NEXT", h: "Your genre track", p: "Carry the same split into specializations." },
      ],
    },
    code: `// bomb:
bullets = withoutNear(bullets, x, y, r) // play
g.fx.FlashScreen(1, 200, 220, 255)
g.fx.Ring(x, y, 3, white)
g.fx.Burst(x, y, 80, 6, white, true)   // fx only`,
  },
];
