/**
 * Advanced Visual Effects Lab lessons (12 chapters).
 * Each remakes a core LEVEL 01–12 game with a separate FX layer so learners
 * notice that particles must not share the gameplay entity lists.
 */
export const advancedLessons = [
  {
    slug: "vfx-fx-tap",
    step: "A01",
    tier: "advanced",
    core: "tap-target",
    stars: "★★★★☆",
    labKind: "fx-split",
    concept: { ja: "イベントでFXを発生", en: "Spawn FX on events" },
    hubDesc: {
      ja: "LEVEL 01 タッチゲーム。当たった瞬間に粒とリングを出す。",
      en: "LEVEL 01 tap game—burst and ring on every hit.",
    },
    ja: {
      navConcept: "タッチ＋キラキラ",
      title: "タッチにエフェクトを足す",
      lead: "応用編の入口。LEVEL 01 の「光る丸をタッチ」に、命中の粒とリングを足します。ゲームの点数処理のあとに、fx.Update() を別呼び出ししていることに注目してください。",
      deepEyebrow: "DEEP DIVE / EVENT → FX",
      deepH: "当たった瞬間、<br>誰が粒を出す？",
      deepLead: "スコアやタイマーは updatePlay()。キラキラは g.fx.Burst() で別の箱へ入れる。同じ Update の中でも、役割を分けると後で迷子になりません。まだ一つのファイルでも、フィールドを分ける練習です。",
      concepts: [
        { h: "プレイ", p: "得点・時間・的の位置。", code: "updatePlay()" },
        { h: "演出", p: "粒・リング・フラッシュ。", code: "fx.Burst / Ring" },
        { h: "順番", p: "ルールを進めてから演出を進める。", code: "play → fx.Update" },
      ],
      lab: {
        eyebrow: "TRY IT / 箱を分ける",
        title: "点数は play、キラキラは fx",
        p: "「命中！」でスコアと粒が同時に増える。「fx だけ 1F」で粒の寿命だけ減る。スコアは動かない——これが分離。",
      },
      play: {
        title: "タッチして花火を見よう",
        p: "丸をタップ。命中で Burst＋Ring。スライダーで演出の強さを変え、上の LIVE GO で play と fx が分かれているのを確認。",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "Update を二段に分ける", p: "ルールと演出を同じ関数にベタ書きしない。" },
      whys: [
        { eyebrow: "WHY SEPARATE?", h: "消えるものが違う", p: "的は次の場所へ移る。粒は寿命で消える。寿命の考え方が違う。" },
        { eyebrow: "WHY ON EVENT?", h: "きっかけはプレイ側", p: "「当たった」はゲームの判定。そこから演出を呼ぶ。" },
        { eyebrow: "TRY NEXT", h: "メーターへ", p: "次は状態（Perfect/OK）ごとに違う演出を出します。" },
      ],
    },
    en: {
      navConcept: "Tap + sparkles",
      title: "Add Effects to Tap",
      lead: "Advanced track starts here. Take LEVEL 01’s tap-the-circle and add hit bursts and rings. Notice updatePlay() then a separate fx.Update().",
      deepEyebrow: "DEEP DIVE / EVENT → FX",
      deepH: "Who spawns sparks<br>on a hit?",
      deepLead: "Score and timer live in updatePlay(). Sparks go into g.fx.Burst(). Same Update frame—different jobs. Even in one file, split the fields.",
      concepts: [
        { h: "Play", p: "Score, time, target position.", code: "updatePlay()" },
        { h: "FX", p: "Particles, rings, flashes.", code: "fx.Burst / Ring" },
        { h: "Order", p: "Advance rules, then effects.", code: "play → fx.Update" },
      ],
      lab: {
        eyebrow: "TRY IT / TWO BOXES",
        title: "Score is play, sparks are fx",
        p: "Hit! grows score and sparks together. Tick FX ages sparks only—score stays put. That’s the split.",
      },
      play: {
        title: "Tap for fireworks",
        p: "Tap the circle for Burst+Ring. Drag intensity and read LIVE GO: play vs fx.",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "Split Update in two", p: "Don’t paste rules and FX into one blob." },
      whys: [
        { eyebrow: "WHY SEPARATE?", h: "Different lifetimes", p: "Targets move; sparks die by life timer." },
        { eyebrow: "WHY ON EVENT?", h: "Play owns the trigger", p: "“Hit” is a game test—then call FX." },
        { eyebrow: "TRY NEXT", h: "Meter", p: "Next: different FX per Perfect/OK/Miss." },
      ],
    },
    code: `func (g *Game) Update() error {
  g.updatePlay() // score, timer, targets
  g.fx.Update()  // particles live here
  return nil
}
// on hit:
g.fx.Burst(x, y, n, speed, tint, true)
g.fx.Ring(x, y, 1.2, tint)`,
  },
  {
    slug: "vfx-fx-meter",
    step: "A02",
    tier: "advanced",
    core: "timing-meter",
    stars: "★★★★☆",
    labKind: "fx-split",
    concept: { ja: "状態ごとに違う演出", en: "FX per game state" },
    hubDesc: {
      ja: "LEVEL 02 メーター。Perfect / OK / Miss で演出を変える。",
      en: "LEVEL 02 meter—different FX for Perfect / OK / Miss.",
    },
    ja: {
      navConcept: "判定→演出の対応表",
      title: "タイミング判定に演出を振る",
      lead: "止めた位置の「状態」が Perfect なら大花火、Miss なら小さな煙。ゲーム状態の enum/分岐と、演出の強さを対応させる練習です。",
      deepEyebrow: "DEEP DIVE / STATE → FX",
      deepH: "同じタップでも、<br>演出が違う理由",
      deepLead: "分岐はプレイ側（ゾーン判定）。そこから fx.Burst の個数や色を変える。演出側に「ゾーン判定」を持ち込まないのがポイントです。",
      concepts: [
        { h: "状態", p: "Perfect / OK / Miss。", code: "result := grade()" },
        { h: "対応", p: "状態ごとに Burst の強さ。", code: "switch result" },
        { h: "寿命", p: "演出は frames で消える。", code: "fx.Update()" },
      ],
      lab: {
        eyebrow: "TRY IT / TWO BOXES",
        title: "判定と花火は別物",
        p: "判定ロジックと粒の寿命は、別の箱の話です。",
      },
      play: {
        title: "ゾーンで止めよう",
        p: "メーターをタップで停止。判定に応じて演出が変わる。LIVE GO の play / fx を見比べよう。",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "switch はプレイ側", p: "演出関数は「どれだけ派手か」だけ受け取る。" },
      whys: [
        { eyebrow: "WHY SWITCH IN PLAY?", h: "ルールの場所", p: "ゲームの正しさはプレイ層。見た目はその結果。" },
        { eyebrow: "WHY TUNE FX?", h: "フィードバック", p: "上手いほど派手、だと上達が気持ちいい。" },
        { eyebrow: "TRY NEXT", h: "キャッチへ", p: "次は「星の配列」と「粒の配列」が同時に増えます。" },
      ],
    },
    en: {
      navConcept: "Grade → FX table",
      title: "FX for Timing Grades",
      lead: "Stop in the zone: Perfect gets a big burst, Miss a gray puff. Map game states to FX strength.",
      deepEyebrow: "DEEP DIVE / STATE → FX",
      deepH: "Same tap,<br>different fireworks",
      deepLead: "Branching stays in play (zone test). Then pick Burst count/color. Don’t put zone math inside the FX system.",
      concepts: [
        { h: "State", p: "Perfect / OK / Miss.", code: "result := grade()" },
        { h: "Map", p: "Strength per result.", code: "switch result" },
        { h: "Life", p: "FX dies on a timer.", code: "fx.Update()" },
      ],
      lab: {
        eyebrow: "TRY IT / TWO BOXES",
        title: "Grade ≠ sparks",
        p: "Judging logic and particle lifetime are different boxes.",
      },
      play: {
        title: "Stop in the zone",
        p: "Tap to freeze the meter. FX changes with the grade—watch play vs fx in LIVE GO.",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "switch stays in play", p: "FX only receives “how flashy”." },
      whys: [
        { eyebrow: "WHY SWITCH IN PLAY?", h: "Rules live there", p: "Correctness is play; looks are the result." },
        { eyebrow: "WHY TUNE FX?", h: "Feedback", p: "Better play → bigger fireworks feels good." },
        { eyebrow: "TRY NEXT", h: "Catch", p: "Next: stars[] and fx.Parts[] both grow." },
      ],
    },
    code: `switch grade(pos) {
case perfect: g.fx.Burst(...); g.fx.FlashScreen(0.8, ...)
case ok:      g.fx.Burst(...smaller)
case miss:    g.fx.Burst(...gray)
}
g.fx.Update()`,
  },
  {
    slug: "vfx-fx-catch",
    step: "A03",
    tier: "advanced",
    core: "catch-stars",
    stars: "★★★★☆",
    labKind: "fx-split",
    concept: { ja: "エンティティとFXを混ぜない", en: "Don’t mix entities & FX" },
    hubDesc: {
      ja: "LEVEL 03 キャッチ。星の配列と粒の配列を分けたまま増やす。",
      en: "LEVEL 03 catch—grow stars[] and fx.Parts[] separately.",
    },
    ja: {
      navConcept: "配列は分けろ",
      title: "落ちものに演出を足すと気づくこと",
      lead: "星が増え、キャッチのたびに粒も増える。もし全部を一つの []any に入れたら、Update で「これは星？粒？」と毎回聞くことになります。ここで分けたくなるはずです。",
      deepEyebrow: "DEEP DIVE / TWO SLICES",
      deepH: "なぜ星と粒を<br>同じスライスにしない？",
      deepLead: "星はカゴ判定と床判定がある。粒は寿命だけ。型も更新も違う。デモは stars と fx.Parts を並べて見せ、「混ぜると地獄」を先に体感させます。",
      concepts: [
        { h: "星", p: "落ちて、当たって、消える遊び。", code: "stars []star" },
        { h: "粒", p: "出して、動いて、寿命で消える。", code: "fx.Parts" },
        { h: "気づき", p: "役割が違うなら箱も分ける。", code: "type Game { stars; fx }" },
      ],
      lab: {
        eyebrow: "TRY IT / TWO BOXES",
        title: "混ぜたら迷子",
        p: "左のリスト＝星、右＝粒。一本化しない。",
      },
      play: {
        title: "星を受け止めよう",
        p: "カゴを動かして星をキャッチ。花火は fx 側。LIVE GO の stars と fx.parts の数を見よう。",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "二つのスライス", p: "Gameplay entities ≠ ephemeral FX." },
      whys: [
        { eyebrow: "WHY TWO SLICES?", h: "問いが違う", p: "「カゴに入った？」と「寿命ゼロ？」は別の if。" },
        { eyebrow: "WHY NOW?", h: "数が増える前に", p: "オブジェクトが増えるほど、分け得が大きくなる。" },
        { eyebrow: "TRY NEXT", h: "フライトへ", p: "次は重力プレイと羽ばたき演出を分離します。" },
      ],
    },
    en: {
      navConcept: "Split the slices",
      title: "Catch Stars—Then Notice the Split",
      lead: "Stars fall; catches spawn sparks. One []any would force “is this a star or a spark?” every frame. You’ll want two boxes.",
      deepEyebrow: "DEEP DIVE / TWO SLICES",
      deepH: "Why not one slice<br>for stars and sparks?",
      deepLead: "Stars need basket/floor tests. Sparks only need life. Different types, different updates. The demo shows stars vs fx.Parts side by side.",
      concepts: [
        { h: "Stars", p: "Fall, collide, score.", code: "stars []star" },
        { h: "Sparks", p: "Spawn, drift, die.", code: "fx.Parts" },
        { h: "Insight", p: "Different jobs → different fields.", code: "type Game { stars; fx }" },
      ],
      lab: {
        eyebrow: "TRY IT / TWO BOXES",
        title: "Mixing gets messy",
        p: "Left list = stars, right = sparks. Keep them apart.",
      },
      play: {
        title: "Catch the stars",
        p: "Move the basket. Fireworks are fx. Watch stars vs fx.parts counts in LIVE GO.",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "Two slices", p: "Gameplay entities ≠ ephemeral FX." },
      whys: [
        { eyebrow: "WHY TWO SLICES?", h: "Different questions", p: "“In basket?” vs “life==0?” are different ifs." },
        { eyebrow: "WHY NOW?", h: "Before counts explode", p: "The more objects, the more a split pays off." },
        { eyebrow: "TRY NEXT", h: "Flight", p: "Next: gravity play vs flap whoosh FX." },
      ],
    },
    code: `type Game struct {
  stars []star     // play entities
  fx    vfxfx.System // ephemeral FX
}
// NEVER: objects []any{star, spark, ...}`,
  },
  {
    slug: "vfx-fx-flight",
    step: "A04",
    tier: "advanced",
    core: "flappy",
    stars: "★★★★☆",
    labKind: "fx-split",
    concept: { ja: "物理と演出の分離", en: "Physics vs FX" },
    hubDesc: {
      ja: "LEVEL 04 フライト。重力・パイプは play、羽ばたき粒は fx。",
      en: "LEVEL 04 flight—gravity/pipes in play, whoosh in fx.",
    },
    ja: {
      navConcept: "重力≠キラキラ",
      title: "ジャンプに羽ばたき演出を足す",
      lead: "速度と重力はプレイの物理。羽ばたきで出る粒は演出。パイプ通過のリングも fx。ここまで来ると、Game 構造体に fx フィールドがあるのが自然に見えます。",
      deepEyebrow: "DEEP DIVE / PHYSICS ≠ FX",
      deepH: "落ちる力と、<br>見える風は別物",
      deepLead: "vy += gravity は世界のルール。Burst は気持ちよさ。混ぜると「演出を切ったら物理も壊れる」事故が起きます。",
      concepts: [
        { h: "物理", p: "位置・速度・重力・当たり。", code: "vy += gravity" },
        { h: "演出", p: "羽ばたき粒・通過リング。", code: "fx.Burst on flap" },
        { h: "構造", p: "Game が両方を持つが更新は分離。", code: "g.fx System" },
      ],
      lab: {
        eyebrow: "TRY IT / 物理と演出",
        title: "落ちる力と、見える花火は別",
        p: "命中ラボと同じく箱はふたつ。フライトでは「重力・パイプ」が play、「羽ばたき粒・通過花火」が fx。片方だけOFFにできるか想像しよう。",
      },
      play: {
        title: "飛んでパイプをくぐろう",
        p: "タップで羽ばたき。関門通過で大きなリング＋花火。物理は play、キラキラは fx——LIVE GO で確認。",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "flap の中で一回呼ぶだけ", p: "物理更新のループに Burst を埋め込まない。" },
      whys: [
        { eyebrow: "WHY ISOLATE?", h: "切っても遊べる", p: "演出ゼロでもゲームとして成立するのが理想。" },
        { eyebrow: "WHY FIELD?", h: "探しやすい", p: "fx を搜せば演出全部が見つかる。" },
        { eyebrow: "TRY NEXT", h: "Pongへ", p: "次は反射の瞬間に火花、移動中に残像。" },
      ],
    },
    en: {
      navConcept: "Gravity ≠ sparkles",
      title: "Flight with Whoosh FX",
      lead: "Velocity and gravity are play physics. Flap sparks and pass rings are fx. A Game.fx field should feel natural now.",
      deepEyebrow: "DEEP DIVE / PHYSICS ≠ FX",
      deepH: "Falling force vs<br>visible wind",
      deepLead: "vy += gravity is world rules. Burst is juice. Mix them and disabling FX can break physics.",
      concepts: [
        { h: "Physics", p: "Pos, vel, gravity, hits.", code: "vy += gravity" },
        { h: "FX", p: "Whoosh sparks, pass rings.", code: "fx.Burst on flap" },
        { h: "Structure", p: "Game owns both; updates stay split.", code: "g.fx System" },
      ],
      lab: {
        eyebrow: "TRY IT / PHYSICS vs FX",
        title: "Falling force ≠ fireworks",
        p: "Same two boxes as the hit lab. Here gravity/pipes are play; flap sparks and gate fireworks are fx. Imagine turning only one off.",
      },
      play: {
        title: "Fly through pipes",
        p: "Tap to flap. Passing a gate fires big rings + bursts. Physics is play; sparkles are fx—read LIVE GO.",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "Call Burst once on flap", p: "Don’t bury Burst inside the physics integrator." },
      whys: [
        { eyebrow: "WHY ISOLATE?", h: "Playable without juice", p: "Ideally the game still works with FX off." },
        { eyebrow: "WHY FIELD?", h: "Findable", p: "Search fx and you find every effect." },
        { eyebrow: "TRY NEXT", h: "Pong", p: "Next: sparks on bounce, trails while moving." },
      ],
    },
    code: `func (g *Game) flap() {
  g.vy = jumpV          // physics (play)
  g.fx.Burst(g.x, g.y…) // juice (fx)
}`,
  },
  {
    slug: "vfx-fx-pong",
    step: "A05",
    tier: "advanced",
    core: "pong",
    stars: "★★★★☆",
    labKind: "fx-split",
    concept: { ja: "衝突とトレイル", en: "Impact + trails" },
    hubDesc: {
      ja: "LEVEL 05 Pong。反射で火花、移動で残像粒。",
      en: "LEVEL 05 Pong—sparks on hit, trail while moving.",
    },
    ja: {
      navConcept: "衝突と残像",
      title: "ボールに火花とトレイルを足す",
      lead: "パドルや壁に当たったフレームだけ Burst。毎フレーム少しトレイルを足すのも fx。ボール自体の vx,vy は play のままです。",
      deepEyebrow: "DEEP DIVE / IMPACT + TRAIL",
      deepH: "一瞬の火花と、<br>続く残像",
      deepLead: "衝突はイベント、トレイルは継続スポーン。どちらも fx に寄せると、ボール構造体が汚れません。",
      concepts: [
        { h: "衝突", p: "当たったフレームに Burst。", code: "onBounce → Burst" },
        { h: "トレイル", p: "移動中に薄い粒を置く。", code: "each N frames" },
        { h: "本体", p: "ボールは位置と速度だけ。", code: "ball {x,y,vx,vy}" },
      ],
      lab: { eyebrow: "TRY IT / TWO BOXES", title: "ボールはボール", p: "見た目の尾を ball に埋め込まない。" },
      play: {
        title: "反射を派手に",
        p: "パドルを動かしてボールを返せ。衝突火花とトレイルは fx 側。",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "ball は薄く", p: "演出フィールドを ball に追加したくなったら fx へ。" },
      whys: [
        { eyebrow: "WHY TRAIL IN FX?", h: "本体が太る", p: "履歴スライスを ball に持つと物理コードが読みにくい。" },
        { eyebrow: "WHY BURST ON HIT?", h: "イベント駆動", p: "毎フレーム火花を出すと散らかる。" },
        { eyebrow: "TRY NEXT", h: "ブロック崩しへ", p: "壊れるブロックと飛び散る破片の役割分担。" },
      ],
    },
    en: {
      navConcept: "Impact + trail",
      title: "Pong Sparks & Trails",
      lead: "Burst only on paddle/wall hits. Soft trail sparks while moving—still fx. Ball vx,vy stay play.",
      deepEyebrow: "DEEP DIVE / IMPACT + TRAIL",
      deepH: "Instant sparks vs<br>ongoing trail",
      deepLead: "Collisions are events; trails are ongoing spawns. Both belong in fx so ball stays thin.",
      concepts: [
        { h: "Impact", p: "Burst on the hit frame.", code: "onBounce → Burst" },
        { h: "Trail", p: "Soft sparks while moving.", code: "each N frames" },
        { h: "Body", p: "Ball is only pose + velocity.", code: "ball {x,y,vx,vy}" },
      ],
      lab: { eyebrow: "TRY IT / TWO BOXES", title: "Ball stays ball", p: "Don’t embed the visual tail inside ball." },
      play: {
        title: "Flashy bounces",
        p: "Move the paddle. Impact sparks and trails are fx.",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "Keep ball thin", p: "If you want an FX field on ball—put it in fx instead." },
      whys: [
        { eyebrow: "WHY TRAIL IN FX?", h: "Body bloat", p: "History slices on ball muddy physics code." },
        { eyebrow: "WHY BURST ON HIT?", h: "Event-driven", p: "Sparks every frame become noise." },
        { eyebrow: "TRY NEXT", h: "Breakout", p: "Bricks vs flying shards." },
      ],
    },
    code: `if bounced {
  g.fx.Burst(ball.x, ball.y, …)
}
if frame%3 == 0 {
  g.fx.Burst(ball.x, ball.y, 1, …) // trail crumb
}`,
  },
  {
    slug: "vfx-fx-breakout",
    step: "A06",
    tier: "advanced",
    core: "breakout",
    stars: "★★★★★",
    labKind: "fx-split",
    concept: { ja: "破壊＝エンティティ削除＋FX", en: "Destroy = remove + FX" },
    hubDesc: {
      ja: "LEVEL 06 ブロック崩し。brick を消し、破片は fx へ。",
      en: "LEVEL 06 breakout—delete brick, shards go to fx.",
    },
    ja: {
      navConcept: "消す／飛ばす",
      title: "ブロック破壊を二段にする",
      lead: "壊れたブロックは配列から削除（プレイ）。飛び散る破片・炎は fx。『壊れたオブジェクトを破片として残す』のを play 側でやると、当たり判定が残骸に引っかかります。",
      deepEyebrow: "DEEP DIVE / DESTROY",
      deepH: "消えるものと、<br>飛び散るもの",
      deepLead: "プレイ空間から消す＝もう当たない。演出空間に出す＝見た目だけ。この二段がアーケード感の正体です。",
      concepts: [
        { h: "削除", p: "bricks から外す。", code: "bricks = alive" },
        { h: "破片", p: "同じ座標で Burst/Flame。", code: "fx.FlameBurst" },
        { h: "分離", p: "残骸に当たり判定を残さない。", code: "FX ≠ collider" },
      ],
      lab: { eyebrow: "TRY IT / TWO BOXES", title: "当たり判定は残すな", p: "見た目の破片をクリックできない／当たらない。" },
      play: {
        title: "ブロックを砕こう",
        p: "ボールでブロックを壊す。破片は fx。bricks の数と fx.parts を見比べよう。",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "delete してから spawn", p: "順序を逆にすると一瞬二重になる。" },
      whys: [
        { eyebrow: "WHY DELETE FIRST?", h: "ルールの真実", p: "もうブロックは無い、が先。花火は後。" },
        { eyebrow: "WHY NOT DEBRIS COLLIDE?", h: "理不尽防止", p: "破片に当たって死ぬとプレイヤーは怒る。" },
        { eyebrow: "TRY NEXT", h: "Snakeへ", p: "体のセル配列とキラキラを分離。" },
      ],
    },
    en: {
      navConcept: "Remove / spray",
      title: "Break Bricks in Two Steps",
      lead: "Remove the brick from play; shards and flames go to fx. Keeping debris as colliders makes unfair hits.",
      deepEyebrow: "DEEP DIVE / DESTROY",
      deepH: "What vanishes vs<br>what sprays",
      deepLead: "Leave play space = no more hits. Enter FX space = looks only. That two-step is arcade juice.",
      concepts: [
        { h: "Delete", p: "Drop from bricks.", code: "bricks = alive" },
        { h: "Shards", p: "Burst/Flame at that spot.", code: "fx.FlameBurst" },
        { h: "Split", p: "FX must not collide.", code: "FX ≠ collider" },
      ],
      lab: { eyebrow: "TRY IT / TWO BOXES", title: "No hitbox on debris", p: "Pretty shards shouldn’t be clickable/collidable." },
      play: {
        title: "Smash bricks",
        p: "Break blocks with the ball. Shards are fx—compare bricks vs fx.parts.",
      },
      codeHead: { eyebrow: "IN THE DEMO", h: "Delete, then spawn", p: "Reverse the order and you double for a frame." },
      whys: [
        { eyebrow: "WHY DELETE FIRST?", h: "Truth of rules", p: "The brick is gone—then fireworks." },
        { eyebrow: "WHY NOT DEBRIS COLLIDE?", h: "Fairness", p: "Dying to confetti angers players." },
        { eyebrow: "TRY NEXT", h: "Snake", p: "Body cells vs sparkles." },
      ],
    },
    code: `// brick broken:
bricks = remove(bricks, i)     // play
g.fx.FlameBurst(x, y, 12)      // fx only`,
  },
];
