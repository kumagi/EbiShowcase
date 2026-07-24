/**
 * The advanced VFX course is deliberately architecture-first. Each chapter
 * revisits one core game, adds a visibly richer motion language, and teaches a
 * different deterministic presentation model. Keep these examples aligned
 * with internal/vfxmotion and games/tracks/visual-effects/vfx-fx-*.
 */

function localized(step, copy) {
  return {
    navConcept: copy.nav,
    title: copy.title,
    lead: copy.lead,
    deepEyebrow: `DEEP DIVE / ${copy.kicker}`,
    deepH: copy.question,
    deepLead: copy.explain,
    concepts: copy.concepts,
    lab: {
      eyebrow: "TRACE IT / RULE → MOTION → PIXELS",
      title: copy.lab.title,
      p: copy.lab.p,
    },
    play: copy.play,
    codeHead: {
      eyebrow: "IN THE DEMO",
      h: copy.codeHead.title,
      p: copy.codeHead.p,
    },
    whys: [
      { eyebrow: "TEST THE CONTRACT", h: copy.test.h, p: copy.test.p },
      { eyebrow: "FAILURE MODE", h: copy.failure.h, p: copy.failure.p },
      { eyebrow: "DESIGN PAYOFF", h: copy.payoff.h, p: copy.payoff.p },
    ],
    architecture: {
      eyebrow: `PATTERN ${step} / ${copy.kicker}`,
      ...copy.architecture,
    },
  };
}

function lesson(meta) {
  return {
    slug: meta.slug,
    step: meta.step,
    tier: "advanced",
    core: meta.core,
    stars: meta.stars || "★★★★★",
    labKind: meta.labKind || "fx-split",
    concept: meta.concept,
    hubDesc: meta.hubDesc,
    ja: localized(meta.step, meta.ja),
    en: localized(meta.step, meta.en),
    code: meta.code,
  };
}

export const allAdvancedLessons = [
  lesson({
    slug: "vfx-fx-tap",
    step: "A01",
    core: "tap-target",
    stars: "★★★★☆",
    concept: { ja: "判定結果を値にして一度だけ渡す", en: "Pass one immutable result" },
    hubDesc: {
      ja: "タッチ判定を TapResult にし、命中の潰れ・残像・衝撃波を安全に重ねる。",
      en: "Turn a tap into TapResult, then layer squash, echo, and shockwave safely.",
    },
    ja: {
      nav: "イベント値とワンショット",
      title: "タッチ結果を演出イベントにする",
      lead: "最初の章では「当たった瞬間」を関数呼び出しの勢いで済ませません。判定を TapResult という小さな値に確定し、得点と演出が同じ事実を一度だけ読む構造へ変えます。命中した的の残像、潰れ、二重の衝撃波、粒子を足しても判定式は一行も増えません。",
      kicker: "IMMUTABLE RESULT",
      question: "命中判定を、<br>どう一度だけ演出へ渡す？",
      explain: "ResolveTap は座標と半径だけから Hit / Miss・表示位置・得点差分を返す純粋関数です。updatePlay は先に得点を確定し、その結果を演出側が消費します。Draw は保存された hitEcho を読むだけなので、残像が何フレーム続いても二重加点しません。",
      concepts: [
        { h: "事実を値にする", p: "入力を Hit / Miss と表示位置へ確定する。", code: "result := ResolveTap(...)" },
        { h: "一度だけ適用", p: "ScoreDelta は Update で一度だけ加える。", code: "score += result.ScoreDelta" },
        { h: "表示用の時間", p: "命中位置の残像は Tween で寿命を持つ。", code: "hitEcho.Advance()" },
      ],
      lab: { title: "同じ結果から数字と花火を分岐させる", p: "命中結果を固定したまま、残像の長さや粒子数だけを変えてください。得点が変わったら境界を越えています。" },
      play: { title: "的を連続で当て、残像を追おう", p: "的は即座に次の場所へ移りますが、当たった場所の輪郭は少し残ります。ゲーム上の的と表示上の命中エコーが別状態であることを観察してください。" },
      codeHead: { title: "判定関数は演出を知らない", p: "純粋な TapResult を境界にし、Shockwave や Burst は結果の消費側で呼びます。" },
      test: { h: "境界・内側・外側を表にする", p: "中心、半径ちょうど、半径の一歩外をパラメータ化し、Outcome と ScoreDelta を比較します。画像もパーティクルも不要です。" },
      failure: { h: "粒子の存在で得点しない", p: "「リングがまだあるから命中中」のように表示状態をルールへ戻すと、フレーム数で得点が変わります。" },
      payoff: { h: "サウンドや振動も同じ結果を読める", p: "TapResult は画面だけのものではありません。音・振動・実績通知を同じ一回の事実へ安全に接続できます。" },
      architecture: {
        name: "Immutable Result",
        problem: "入力判定の中へ得点・粒子・音を直書きすると、演出を増やすたびに判定が壊れます。まず「何が起きたか」を不変の値として確定します。",
        flow: ["tap input", "ResolveTap", "TapResult", "score + hitEcho", "Draw / FX"],
        stateTitle: "hitEcho は命中位置と Tween だけ",
        state: "的そのものは次へ移動して構いません。消えた命中瞬間を見せるため、表示層は結果の座標と14フレームの時計だけを所有します。",
        testTitle: "テストするのは結果の事実",
        tests: ["半径境界で Hit になる", "Miss の表示位置はタップ位置になる", "ScoreDelta は一回分だけ返る"],
      },
    },
    en: {
      nav: "Result events & one-shots",
      title: "Turn a Tap Result into Motion",
      lead: "Do not hide “the instant of a hit” inside a pile of side effects. Resolve the input into a small TapResult value. Score and presentation consume the same fact once. Target echo, squash, a double shockwave, and particles can grow without growing the hit test.",
      kicker: "IMMUTABLE RESULT",
      question: "How does a hit cross<br>into presentation once?",
      explain: "ResolveTap is a pure function of coordinates and radius. It returns Hit / Miss, a presentation position, and a score delta. updatePlay commits score first; presentation consumes the result. Draw only reads hitEcho, so a fourteen-frame echo can never award fourteen hits.",
      concepts: [
        { h: "Make a fact", p: "Resolve input to outcome and position.", code: "result := ResolveTap(...)" },
        { h: "Apply it once", p: "ScoreDelta is committed in Update once.", code: "score += result.ScoreDelta" },
        { h: "Own visual time", p: "The hit echo gets an independent Tween.", code: "hitEcho.Advance()" },
      ],
      lab: { title: "Fork numbers and fireworks from one result", p: "Keep TapResult fixed while changing echo duration and particle count. If score changes, presentation crossed the boundary." },
      play: { title: "Chain hits and follow their echoes", p: "The target moves immediately; the old hit outline lingers. Compare the gameplay target with the presentation-only echo." },
      codeHead: { title: "The judge knows no effects", p: "Use pure TapResult as the seam; Shockwave and Burst belong to its consumer." },
      test: { h: "Table-test inside, edge, and outside", p: "Compare Outcome and ScoreDelta at the center, exactly on the radius, and just outside. No image or particle system is needed." },
      failure: { h: "Never score from living particles", p: "If “a ring still exists” means the hit is active, effect duration changes the score." },
      payoff: { h: "Sound and haptics can read the same fact", p: "TapResult is not screen-specific. Audio, vibration, and achievements can subscribe to the same one-shot result." },
      architecture: {
        name: "Immutable Result",
        problem: "Inlining score, particles, and sound in the hit test makes every new flourish a rule change. First freeze what happened into a value.",
        flow: ["tap input", "ResolveTap", "TapResult", "score + hitEcho", "Draw / FX"],
        stateTitle: "hitEcho owns only position and a Tween",
        state: "The target may move on immediately. Presentation keeps the resolved position and a fourteen-frame clock to show the vanished instant.",
        testTitle: "Test the fact, not the fireworks",
        tests: ["the radius boundary is a hit", "a miss retains the tap position", "ScoreDelta represents exactly one hit"],
      },
    },
    code: `result := vfxmotion.ResolveTap(x, y, g.cx, g.cy, g.r)
if result.Outcome == vfxmotion.TapHit {
    g.score += result.ScoreDelta       // commit the rule once
    g.hitEcho = vfxmotion.NewTween(14) // presentation clock
    g.fx.Shockwave(result.X, result.Y, .75, inner, outer)
}`,
  }),

  lesson({
    slug: "vfx-fx-meter",
    step: "A02",
    core: "timing-meter",
    labKind: "fx-meter-grade",
    concept: { ja: "判定enumを演出レシピへ写す", en: "Map a grade enum to a recipe" },
    hubDesc: {
      ja: "Perfect / OK / Missを型付きレシピへ変換し、停止・バナー・花火を一貫させる。",
      en: "Map Perfect / OK / Miss to typed recipes for freeze, banner, and fireworks.",
    },
    ja: {
      nav: "型付き演出レシピ",
      title: "タイミング判定を演出レシピへ写す",
      lead: "Perfectの得点と花火を別々のswitchにすると、調整の途中で食い違います。Gradeを一度求め、Score・BurstCount・Flash・Freezeを持つデータだけのRecipeへ変換します。大成功ほど停止感、バナーの膨張、衝撃波が強くなります。",
      kicker: "TYPED RECIPE",
      question: "判定ごとの派手さを、<br>どう一箇所で管理する？",
      explain: "距離計算は JudgeMeter、見せ方の対応表は RecipeForGrade に分けます。FXシステムに「中心から10pxならPerfect」と教えてはいけません。演出はGradeだけを受け取り、数値の表を読むため、難易度調整と見た目調整を別々にレビューできます。",
      concepts: [
        { h: "ルールのenum", p: "位置を Grade に確定する。", code: "grade := JudgeMeter(...)" },
        { h: "表示のレシピ", p: "得点・粒数・停止時間を一組にする。", code: "recipe := RecipeForGrade(grade)" },
        { h: "一つのフィードバック時計", p: "バナーとマーカーの伸縮を同期する。", code: "feedback Tween" },
      ],
      lab: { title: "Gradeを固定してレシピだけ調整する", p: "PerfectのFreezeやBurstCountを変え、判定範囲が変化しないことを確かめます。対応表は視覚調整の安全な操作面です。" },
      play: { title: "三つの判定をわざと出して比べる", p: "Perfectではマーカーが大きく止まり、OKは中程度、Missは重い土煙だけ。点数と見た目が同じGradeを共有します。" },
      codeHead: { title: "switchではなくレシピを読む", p: "JudgeとRecipeを別の純粋関数にすると、境界値と演出対応を独立して表形式で試せます。" },
      test: { h: "Grade×Recipeの表を比較する", p: "marker位置から期待Gradeを確認し、そのGradeのScore・Label・BurstCountを同じテーブルで検証します。" },
      failure: { h: "FX側で距離を再計算しない", p: "表示側が独自にPerfectを判定すると、点数はOKなのに画面はPerfectという矛盾が生まれます。" },
      payoff: { h: "デザイナーが数値表だけ調整できる", p: "レシピがデータなら、後でJSON化しても判定アルゴリズムに触れずに手触りを調整できます。" },
      architecture: {
        name: "Typed Presentation Recipe",
        problem: "判定・得点・色・粒数を複数のswitchへ散らすと、同じPerfectが場所ごとに別の意味になります。",
        flow: ["marker position", "Grade", "GradeRecipe", "feedback Tween", "banner + FX"],
        stateTitle: "feedback はGradeと一個の時計を持つ",
        state: "マーカー、バナー、フラッシュを別々のタイマーでずらさず、同じProgressから大きさと明度を導出します。",
        testTitle: "境界と対応表を別々に保証",
        tests: ["中心±10はPerfect", "Perfectは100点と最大Burst", "Missは得点0・Freeze 0"],
      },
    },
    en: {
      nav: "Typed FX recipe",
      title: "Map Timing Grades to Recipes",
      lead: "Separate switches for Perfect score and Perfect fireworks will drift apart. Resolve Grade once, then map it to a data-only recipe containing Score, BurstCount, Flash, and Freeze. Better timing now produces stronger pause, banner expansion, and shockwave from one source.",
      kicker: "TYPED RECIPE",
      question: "How do all grade effects<br>stay in one place?",
      explain: "JudgeMeter owns distance math; RecipeForGrade owns presentation mapping. Never teach the FX system that ten pixels means Perfect. It receives Grade and reads numbers, so difficulty and visual tuning can be reviewed separately.",
      concepts: [
        { h: "Rule enum", p: "Resolve a position into Grade.", code: "grade := JudgeMeter(...)" },
        { h: "Presentation recipe", p: "Bundle score, count, and freeze.", code: "recipe := RecipeForGrade(grade)" },
        { h: "One feedback clock", p: "Synchronize marker and banner motion.", code: "feedback Tween" },
      ],
      lab: { title: "Freeze Grade; tune only its recipe", p: "Change Perfect Freeze or BurstCount and verify that the judgement window never moves. The recipe is a safe tuning surface." },
      play: { title: "Deliberately produce all three grades", p: "Perfect holds and swells, OK is moderate, and Miss makes only a heavy puff. Score and visuals share one Grade." },
      codeHead: { title: "Read a recipe instead of branching everywhere", p: "Separate pure Judge and Recipe functions make boundaries and presentation mapping independently table-testable." },
      test: { h: "Test a Grade × Recipe table", p: "Resolve marker positions, then compare each recipe’s Score, Label, and BurstCount in the same table." },
      failure: { h: "Do not re-judge distance in FX", p: "A second Perfect calculation can show Perfect while gameplay awards only OK." },
      payoff: { h: "Designers can tune a number table", p: "A data recipe can later move to JSON without exposing the judgement algorithm." },
      architecture: {
        name: "Typed Presentation Recipe",
        problem: "Scattering grade, score, color, and count across switches gives Perfect a different meaning in every file.",
        flow: ["marker position", "Grade", "GradeRecipe", "feedback Tween", "banner + FX"],
        stateTitle: "feedback owns Grade and one clock",
        state: "Marker, banner, and flash derive from the same Progress instead of drifting on separate timers.",
        testTitle: "Guarantee boundaries and mapping separately",
        tests: ["center ±10 is Perfect", "Perfect maps to 100 and the largest burst", "Miss maps to zero score and freeze"],
      },
    },
    code: `grade := vfxmotion.JudgeMeter(g.marker, center)
recipe := vfxmotion.RecipeForGrade(grade)
g.score += recipe.Score
g.feedback = vfxmotion.NewTween(22 + recipe.Freeze)
g.spawnGradeFX(grade, recipe)`,
  }),

  lesson({
    slug: "vfx-fx-catch",
    step: "A03",
    core: "catch-stars",
    concept: { ja: "消えた実体を表示代理で見送る", en: "Let a visual proxy outlive an entity" },
    hubDesc: {
      ja: "捕獲した星はルールから即削除し、ID付き表示代理だけをHUDへ飛ばす。",
      en: "Remove a caught star immediately, then fly an ID-backed visual proxy to the HUD.",
    },
    ja: {
      nav: "表示代理のライフサイクル",
      title: "捕まえた星をHUDまで飛ばす",
      lead: "捕獲された星は当たり判定から即座に消すべきですが、画面からも瞬間消去すると手応えがありません。そこで星のIDと捕獲位置を表示代理Proxyへ写し、本物は削除、代理だけを放物線でHUDへ飛ばします。カゴも同じTweenで潰れて戻ります。",
      kicker: "VISUAL PROXY",
      question: "削除済みの星を、<br>どう画面だけに残す？",
      explain: "ゲームのstarsには落下・捕獲・失敗だけを持たせます。捕獲時にNewProxy(id, from, HUD, frames)を作り、starsからは除去します。Proxyは衝突判定を一切持たず、完了したら表示リストから消えます。",
      concepts: [
        { h: "安定したID", p: "どの星の表示かを値で追跡する。", code: "star.id" },
        { h: "即時削除", p: "捕獲済みの星は再判定しない。", code: "stars = next" },
        { h: "表示代理", p: "位置をHUDへ補間し、終われば捨てる。", code: "proxies []Proxy" },
      ],
      lab: { title: "星を消しても軌跡が完走する", p: "捕獲直後にstarsの長さが減り、proxiesが増えることを追います。代理が飛んでいる間に同じ星がもう一度捕獲されてはいけません。" },
      play: { title: "星を受け、左上へ飛ぶ光を追う", p: "星はカゴで弾け、光の代理が弧を描いてHUDへ移動します。カゴの潰れと衝撃波も捕獲イベント一回から始まります。" },
      codeHead: { title: "実体を延命せず、表示情報だけコピー", p: "アニメーションのためにcaughtフラグを星へ足し続けず、別の短命モデルへ必要最小限を写します。" },
      test: { h: "開始点・終点・完了を試す", p: "Proxyの0Fは捕獲位置、最終FはHUD位置、指定フレーム後にDoneとなることを画像なしで確認します。" },
      failure: { h: "捕獲済み実体を当たり判定へ残さない", p: "見せるために星を残すと、カゴの中で毎フレーム再捕獲される事故が起きます。" },
      payoff: { h: "コイン・経験値・アイテム取得に再利用", p: "拾った物がUIへ吸い込まれる表現は、ゲーム実体と表示代理の分離で一貫して作れます。" },
      architecture: {
        name: "Visual Proxy Lifecycle",
        problem: "取得物をアニメーションのためルール側へ残すと、衝突や保存データに「取得済みだが存在する物」が混ざります。",
        flow: ["star ID", "catch event", "remove star", "spawn Proxy", "fly to HUD"],
        stateTitle: "Proxyはfrom/toとTweenだけ",
        state: "代理には速度判定・得点・当たり判定を持たせません。表示に必要なIDと二点、時計だけで完結させます。",
        testTitle: "実体と代理の寿命を別に試す",
        tests: ["捕獲後starsから一回で消える", "Proxyは弧を描いてHUDへ着く", "Done後はproxiesから除去される"],
      },
    },
    en: {
      nav: "Visual-proxy lifecycle",
      title: "Fly a Caught Star to the HUD",
      lead: "A caught star should leave collision immediately, but vanishing from the screen instantly feels weak. Copy its ID and catch position into a presentation Proxy, delete the real star, and fly only the proxy along an arc to the HUD. The basket squashes on the same catch event.",
      kicker: "VISUAL PROXY",
      question: "How can a deleted star<br>remain only on screen?",
      explain: "Gameplay stars own fall, catch, and miss. A catch creates NewProxy(id, from, HUD, frames), then removes the star. Proxy has no collision and leaves the presentation list when its tween finishes.",
      concepts: [
        { h: "Stable ID", p: "Track which visual was captured.", code: "star.id" },
        { h: "Immediate removal", p: "Never collide a caught star again.", code: "stars = next" },
        { h: "Visual proxy", p: "Interpolate to HUD, then discard.", code: "proxies []Proxy" },
      ],
      lab: { title: "Delete the star while its trip completes", p: "Immediately after catch, stars shrinks and proxies grows. The same star cannot score again while its proxy is flying." },
      play: { title: "Catch a star and follow the light left", p: "The star bursts at the basket and a visual proxy arcs into the HUD. Basket squash and shockwave share the one catch event." },
      codeHead: { title: "Copy display facts instead of extending the entity", p: "Do not keep adding caught flags to gameplay stars; copy the minimum into a short-lived model." },
      test: { h: "Test start, destination, and completion", p: "At frame zero Proxy is at catch position, at the last frame it reaches HUD, then Done is true—all without images." },
      failure: { h: "Never leave a caught collider alive", p: "Keeping the star for animation lets the basket catch it again every frame." },
      payoff: { h: "Reuse for coins, XP, and loot", p: "Any pickup that flies into UI benefits from the same entity/proxy split." },
      architecture: {
        name: "Visual Proxy Lifecycle",
        problem: "Keeping pickups in rules for animation pollutes collision and save data with objects that are both collected and present.",
        flow: ["star ID", "catch event", "remove star", "spawn Proxy", "fly to HUD"],
        stateTitle: "Proxy owns only from/to and a Tween",
        state: "It has no scoring, velocity rule, or collider—just an ID, two points, and a clock.",
        testTitle: "Test entity and proxy lifetimes separately",
        tests: ["catch removes one star once", "Proxy arcs to the HUD", "finished proxies are pruned"],
      },
    },
    code: `if caught {
    g.score++                    // gameplay commits now
    g.proxies = append(g.proxies,
        vfxmotion.NewProxy(st.id, st.x, st.y, hudX, hudY, 24))
    continue                     // st is absent from next gameplay slice
}`,
  }),

  lesson({
    slug: "vfx-fx-flight",
    step: "A04",
    core: "flappy",
    concept: { ja: "物理状態から姿勢を導出する", en: "Derive pose from physics" },
    hubDesc: {
      ja: "上昇・滑空・急降下・衝突の姿勢を速度から導き、物理へ逆流させない。",
      en: "Derive flap, glide, dive, and crash poses from velocity without feeding physics back.",
    },
    ja: {
      nav: "物理から導くポーズ",
      title: "飛行物理を読みやすい姿勢へ翻訳する",
      lead: "vyは衝突計算には十分でも、プレイヤーには上昇・失速・落下が読めません。速度、羽ばたきからの経過、衝突済みかをPoseForFlightへ渡し、回転と縦横スケールを導出します。風の流線、羽ばたきの粒、ゲート通過の紙吹雪も加えます。",
      kicker: "DERIVED POSE",
      question: "物理の数字を、<br>どう身体の動きへ翻訳する？",
      explain: "鳥の座標とvyはplayの真実です。FlightPoseとPoseTransformは毎フレームそこから計算できる表示値で、保存も衝突も不要です。DrawCenteredPoseは中心を保ったまま回転・潰れ・伸びだけを適用します。",
      concepts: [
        { h: "物理の真実", p: "重力と羽ばたきはvyを更新する。", code: "birdY += vy" },
        { h: "導出した姿勢", p: "vyと経過からFlap / Glide / Dive。", code: "PoseForFlight(...)" },
        { h: "描画変換", p: "回転とScaleだけを画像へ適用する。", code: "DrawCenteredPose" },
      ],
      lab: { title: "同じvyから必ず同じ姿勢を得る", p: "速度を固定し、羽ばたき直後・滑空・急降下・衝突のPoseを比較します。姿勢を変えてもbirdYの更新式は変えません。" },
      play: { title: "羽ばたきから急降下まで姿勢を見る", p: "タップ直後は縦に伸び、落下が速いほど前へ傾きます。ゲート通過は二重波と紙吹雪、衝突は赤い衝撃波です。" },
      codeHead: { title: "姿勢はDraw直前に導出できる", p: "physicsへanimationFrameを混ぜず、入力となる物理値と短いframesSinceFlapだけを渡します。" },
      test: { h: "速度×経過×衝突を表にする", p: "fresh flap、低速滑空、高速落下、crashedの各入力が期待Poseになることをパラメータ化します。" },
      failure: { h: "傾いた画像の角で衝突しない", p: "見た目の回転を当たり判定へ戻すと、アート調整が難易度変更になります。" },
      payoff: { h: "敵・乗り物・投射物にも使える", p: "物理ベクトルから表示姿勢を導く考え方は、飛行機、魚、ノックバック中のキャラにも共通します。" },
      architecture: {
        name: "Pose Derived from Physics",
        problem: "animationStateを物理と別に手書きすると、上昇しているのに落下絵という不整合が起きます。",
        flow: ["birdY + vy", "framesSinceFlap", "FlightPose", "PoseTransform", "sprite + wind"],
        stateTitle: "保存するのは羽ばたきからの経過だけ",
        state: "GlideやDiveはvyから導出できます。重複状態を保存せず、入力イベントの余韻だけ小さな時計にします。",
        testTitle: "物理を変えず表示関数だけ試す",
        tests: ["羽ばたき7F以内はFlap", "大きい正のvyはDive", "crashedは速度に関係なくCrash"],
      },
    },
    en: {
      nav: "Pose from physics",
      title: "Translate Flight Physics into Pose",
      lead: "vy is enough for collision but not for a player reading rise, stall, and fall. Feed velocity, frames since flap, and crash state to PoseForFlight to derive rotation and squash. Add wind streaks, flap sparks, and gate confetti without touching gravity.",
      kicker: "DERIVED POSE",
      question: "How do physics numbers<br>become body language?",
      explain: "birdY and vy are play truth. FlightPose and PoseTransform are disposable presentation values derived every frame. They need no save state or collider. DrawCenteredPose keeps the center stable while applying rotation and scale.",
      concepts: [
        { h: "Physics truth", p: "Gravity and flap update vy.", code: "birdY += vy" },
        { h: "Derived pose", p: "vy and age select Flap / Glide / Dive.", code: "PoseForFlight(...)" },
        { h: "Draw transform", p: "Apply rotation and scale to the image only.", code: "DrawCenteredPose" },
      ],
      lab: { title: "The same vy always yields the same pose", p: "Hold velocity fixed and compare fresh flap, glide, dive, and crash. Never modify the birdY equation to tune a pose." },
      play: { title: "Read the body from flap to dive", p: "A fresh flap stretches upward; a fast fall tilts forward. Gate passes add double waves and confetti; crashes add a red shockwave." },
      codeHead: { title: "Derive pose just before drawing", p: "Keep animationFrame out of physics; pass physics values and only the short flap-age clock." },
      test: { h: "Table-test velocity × age × crash", p: "Fresh flap, slow glide, fast fall, and crashed inputs should map to their expected Pose." },
      failure: { h: "Do not collide the rotated art", p: "Feeding visual rotation into collision turns an art tweak into a difficulty change." },
      payoff: { h: "Reuse for enemies, vehicles, and projectiles", p: "Deriving display pose from a physics vector works for planes, fish, and knocked-back characters." },
      architecture: {
        name: "Pose Derived from Physics",
        problem: "A separately hand-written animationState can show falling art while physics is rising.",
        flow: ["birdY + vy", "framesSinceFlap", "FlightPose", "PoseTransform", "sprite + wind"],
        stateTitle: "Store only time since the flap",
        state: "Glide and Dive derive from vy. Avoid duplicated state; keep only the short aftertaste of the input event.",
        testTitle: "Test the display function without changing physics",
        tests: ["the first seven frames select Flap", "large positive vy selects Dive", "crashed always selects Crash"],
      },
    },
    code: `g.vy += gravity
g.birdY += g.vy // collision truth

_, pose := vfxmotion.PoseForFlight(g.vy, g.sinceFlap, g.gameOver)
hero.DrawCenteredPose(screen, birdX, g.birdY, 42,
    hero.Pose{Rotation: pose.Rotation, ScaleX: pose.ScaleX, ScaleY: pose.ScaleY})`,
  }),

  lesson({
    slug: "vfx-fx-pong",
    step: "A05",
    core: "pong",
    concept: { ja: "上限つき描画履歴で軌跡を作る", en: "Build trails from bounded history" },
    hubDesc: {
      ja: "ボール本体へ残像を足さず、表示層の固定長履歴で滑らかな軌跡を描く。",
      en: "Keep afterimages out of the ball and draw a smooth trail from bounded presentation history.",
    },
    ja: {
      nav: "固定長の描画履歴",
      title: "ボールの過去を安全な軌跡にする",
      lead: "残像は粒を毎フレーム無制限に出す必要がありません。表示層に直近16点だけのTrailを持ち、古い点から薄く小さく描けば、メモリ上限が明確な滑らかな光跡になります。打球時にはパドルを横へ潰し、接触点へリングと火花を重ねます。",
      kicker: "BOUNDED HISTORY",
      question: "過去の位置を、<br>どこまで覚えるべき？",
      explain: "ボールは現在位置と速度だけで衝突できます。過去位置はゲームルールではなく表示用キャッシュです。Trail.Pushは上限に達したら最古を捨て、Pointsはコピーを返すのでDrawが内部状態を書き換えません。",
      concepts: [
        { h: "現在はplay", p: "衝突はbx/by/vx/vyだけを見る。", code: "ball state" },
        { h: "過去はpresentation", p: "直近16点をリングバッファ相当に保持。", code: "trail.Push(point)" },
        { h: "強さを順番から導出", p: "新しい点ほど大きく明るく描く。", code: "alpha := age / len" },
      ],
      lab: { title: "履歴上限を小さくして軌跡を比べる", p: "3点、8点、16点で軌跡だけが変わり、跳ね返りと得点は同じであることを確認します。" },
      play: { title: "ラリーして光跡とパドルの潰れを見る", p: "ボールの軌跡は連続して見え、接触時だけパドルが横へ広がります。得点後はTrailをClearして別ラリーの残像を混ぜません。" },
      codeHead: { title: "粒子ではなく履歴を描く選択", p: "連続軌跡は固定長サンプルの方が安価で決定的です。衝撃だけ粒子に任せます。" },
      test: { h: "押し出し順と上限を試す", p: "上限3へ4点Pushした時、Pointsが2,3,4だけを返すことを比較します。" },
      failure: { h: "履歴をボールの当たり判定へ混ぜない", p: "過去位置の円まで衝突させると、長い軌跡ほど巨大な当たり判定になります。" },
      payoff: { h: "剣筋・レーザー・リプレイゴーストにも使える", p: "上限つき時系列は、連続線や入力軌跡など粒子より形を保ちたい表示に向きます。" },
      architecture: {
        name: "Bounded Presentation History",
        problem: "毎フレーム粒子を出す軌跡は個数が速度やFPSに依存し、上限と形が読みづらくなります。",
        flow: ["ball position", "Trail.Push", "last 16 points", "age gradient", "glow trail"],
        stateTitle: "Trailは表示用キャッシュ",
        state: "得点や衝突に使わず、サーブ時にClearします。古い点の削除規則が一箇所なのでメモリ使用量が固定されます。",
        testTitle: "履歴の順序と上限だけを保証",
        tests: ["上限まで順番を保つ", "超過時は最古を一つ捨てる", "Clearで次のラリーへ持ち越さない"],
      },
    },
    en: {
      nav: "Bounded draw history",
      title: "Turn Ball History into a Safe Trail",
      lead: "A trail does not need unlimited particles every frame. Keep only the latest sixteen points in presentation and draw old samples smaller and dimmer. Memory is bounded, shape is smooth, and impacts can still add paddle squash, ring, and sparks.",
      kicker: "BOUNDED HISTORY",
      question: "How much of the past<br>should presentation remember?",
      explain: "The ball collides with current position and velocity only. Old positions are a presentation cache. Trail.Push evicts the oldest sample at its cap; Points returns a copy so Draw cannot mutate history.",
      concepts: [
        { h: "Present is play", p: "Collision reads bx/by/vx/vy only.", code: "ball state" },
        { h: "Past is presentation", p: "Keep the latest sixteen points.", code: "trail.Push(point)" },
        { h: "Style from order", p: "Newer points draw larger and brighter.", code: "alpha := age / len" },
      ],
      lab: { title: "Compare short and long bounded trails", p: "Try 3, 8, and 16 samples. Bounce and score remain identical while only the trail changes." },
      play: { title: "Rally through trails and paddle squash", p: "The trail reads as continuous; the paddle widens only on contact. A serve clears history so rallies never share ghosts." },
      codeHead: { title: "Choose history instead of particles", p: "A continuous path is cheaper and deterministic as bounded samples. Reserve particles for impacts." },
      test: { h: "Test eviction order and cap", p: "Pushing four points into a limit-three Trail must return only 2, 3, 4." },
      failure: { h: "Never collide the trail", p: "If old circles collide, a longer trail creates a larger hitbox." },
      payoff: { h: "Reuse for sword arcs, lasers, and replay ghosts", p: "Bounded time series fit visuals that must preserve a shape better than particles do." },
      architecture: {
        name: "Bounded Presentation History",
        problem: "Per-frame particle trails depend on speed and frame rate, obscuring both memory cap and shape.",
        flow: ["ball position", "Trail.Push", "last 16 points", "age gradient", "glow trail"],
        stateTitle: "Trail is a presentation cache",
        state: "It never scores or collides and clears on serve. One eviction rule fixes memory use.",
        testTitle: "Guarantee only order and capacity",
        tests: ["order is preserved up to the limit", "overflow evicts the oldest", "Clear isolates the next rally"],
      },
    },
    code: `g.bx += g.vx * speed
g.by += g.vy * speed
g.trail.Push(vfxmotion.Point{X: g.bx, Y: g.by})

// Draw: old samples are dim; only the current ball collides.
for i, p := range g.trail.Points() { drawTrailSample(p, i) }`,
  }),

  lesson({
    slug: "vfx-fx-breakout",
    step: "A06",
    core: "breakout",
    concept: { ja: "削除前スナップショットで破壊を描く", en: "Animate deletion from a snapshot" },
    hubDesc: {
      ja: "ブロックは即削除し、削除前の形と色をTombstoneへ写して破片にする。",
      en: "Delete the brick now, copy its old shape into a Tombstone, and animate shards afterward.",
    },
    ja: {
      nav: "削除前スナップショット",
      title: "壊れたブロックを破片へ引き継ぐ",
      lead: "ブロックを消してから色や形を読もうとしても、参照先はもうありません。衝突が成立した瞬間にID・矩形・色番号をSnapshotへコピーし、bricksではalive=falseを即確定します。Tombstoneはその静止画を24フレームだけ破片へ崩し、2フレームのヒットストップが衝撃を強調します。",
      kicker: "PRE-DELETE SNAPSHOT",
      question: "削除した物の破片を、<br>何から描けばいい？",
      explain: "削除前に表示に必要な事実だけを値コピーします。Tombstoneは元brickへのポインタを持たないため、配列の詰め替えや次の面の生成後も安全です。ゲーム側ではブロックが消えているので、破片へ再衝突しません。",
      concepts: [
        { h: "先に写す", p: "ID・位置・大きさ・色番号を値コピー。", code: "Snapshot{...}" },
        { h: "ルールは即削除", p: "alive=falseと得点をその場で確定。", code: "brick.alive = false" },
        { h: "墓標を崩す", p: "TombstoneのProgressで四片を広げる。", code: "ghost.Progress()" },
      ],
      lab: { title: "ブロックなしでも破片が完走する", p: "aliveを落とした直後からTombstoneだけを進めます。次の面をbuildBricksしても旧破片の色と位置が変わらないことを確認します。" },
      play: { title: "一個ずつ壊し、破片の余韻を見る", p: "衝突時にボールは短く止まり、ブロックは四片へ割れて落ちます。全消しでは紙吹雪を重ねますが、新面のブロックとは別リストです。" },
      codeHead: { title: "参照ではなく値を保存する", p: "削除対象へのポインタを演出が握ると、再利用された配列要素を誤って描きます。" },
      test: { h: "Snapshotが最後まで不変か試す", p: "Tombstoneを完了までAdvanceし、元のID・矩形・Variantが変化しないことを比較します。" },
      failure: { h: "壊れかけ状態をルールへ増やしすぎない", p: "breaking brickを衝突配列へ残すと、再ヒット・クリア判定・保存形式が複雑になります。" },
      payoff: { h: "敵死亡・宝箱・地形破壊へ拡張", p: "即時削除＋表示スナップショットは、元実体より長い死亡演出すべてに使えます。" },
      architecture: {
        name: "Pre-delete Snapshot",
        problem: "削除後に元実体へ参照して破片を描くと、nil・配列再利用・次ステージのデータ混入が起きます。",
        flow: ["brick hit", "copy Snapshot", "remove brick", "Tombstone Tween", "four shards"],
        stateTitle: "Tombstoneは元実体を参照しない",
        state: "描画に必要な矩形とVariantだけを値で持ち、24F後に捨てます。元のbricks配列とは独立して更新できます。",
        testTitle: "削除と余韻を別の契約にする",
        tests: ["hit直後にalive=false", "Snapshotの値はTween中不変", "Done後にghostsから除去"],
      },
    },
    en: {
      nav: "Pre-delete snapshot",
      title: "Hand a Broken Brick to Its Shards",
      lead: "After deletion there is no safe brick from which to read shape or color. At collision, copy ID, rectangle, and color variant into Snapshot, then commit alive=false immediately. Tombstone breaks that still image into shards for twenty-four frames while a two-frame hit stop sells impact.",
      kicker: "PRE-DELETE SNAPSHOT",
      question: "What draws the shards<br>after the object is gone?",
      explain: "Value-copy only the facts presentation needs before deletion. Tombstone holds no pointer to the brick, so slice compaction and next-level creation are harmless. Rules already consider the brick absent, so shards cannot collide again.",
      concepts: [
        { h: "Copy first", p: "Value-copy ID, rect, and color variant.", code: "Snapshot{...}" },
        { h: "Delete in rules now", p: "Commit alive=false and score.", code: "brick.alive = false" },
        { h: "Break the tombstone", p: "Spread four shards from Progress.", code: "ghost.Progress()" },
      ],
      lab: { title: "Let shards finish without a brick", p: "Advance Tombstone after alive becomes false. Rebuilding bricks must not alter old shard position or color." },
      play: { title: "Break one brick and watch its afterlife", p: "The ball pauses briefly and the brick falls into four pieces. Clear confetti is another presentation list, not the new level." },
      codeHead: { title: "Store values, not a reference", p: "A presentation pointer to a deleted slice element may draw an element later reused for something else." },
      test: { h: "Verify Snapshot stays immutable", p: "Advance Tombstone to completion and compare the original ID, rectangle, and Variant." },
      failure: { h: "Do not grow a breaking rule state unnecessarily", p: "Keeping breaking bricks in collision complicates repeat hits, clear checks, and save data." },
      payoff: { h: "Extend to deaths, chests, and terrain", p: "Immediate deletion plus a visual snapshot fits every death animation that outlives its entity." },
      architecture: {
        name: "Pre-delete Snapshot",
        problem: "Drawing shards through the deleted entity risks nil access, slice reuse, and data from the next stage.",
        flow: ["brick hit", "copy Snapshot", "remove brick", "Tombstone Tween", "four shards"],
        stateTitle: "Tombstone never references the entity",
        state: "It value-owns only the rectangle and Variant needed for drawing, then expires after twenty-four frames.",
        testTitle: "Separate deletion and afterlife contracts",
        tests: ["alive is false immediately", "Snapshot stays unchanged during Tween", "Done tombstones are pruned"],
      },
    },
    code: `snap := vfxmotion.Snapshot{
    ID: i, X: b.x, Y: b.y, W: brickW, H: brickH, Variant: b.row,
}
g.ghosts = append(g.ghosts, vfxmotion.NewTombstone(snap, 24))
b.alive = false // collision truth changes immediately`,
  }),

  lesson({
    slug: "vfx-fx-snake",
    step: "A07",
    core: "snake",
    concept: { ja: "離散グリッドを補間し形を導出する", en: "Interpolate a discrete grid and derive shape" },
    hubDesc: {
      ja: "bodyのマス座標は整数のまま、前回位置との補間と隣接セルから滑らかな蛇を描く。",
      en: "Keep integer body cells while interpolating previous positions and deriving curves from neighbors.",
    },
    ja: {
      nav: "離散状態の補間",
      title: "マス目の蛇を滑らかにつなぐ",
      lead: "Snakeの正解は整数マス列です。そこを小数座標へ変えると自己衝突や餌判定が壊れます。移動tickの直前bodyをvisualFromへコピーし、次tickまでのProgressで各節を補間します。さらに前・現在・次の三セルから直線か角かを導出し、円と太線で切れ目のない身体を描きます。",
      kicker: "DISCRETE INTERPOLATION",
      question: "整数の盤面を保ったまま、<br>どう滑らかに見せる？",
      explain: "ルール更新は従来どおり一マス跳びます。表示だけがold→newを9フレームで移動します。SegmentPoseは隣接関係からStraight / Cornerを返す純粋関数なので、spriteIndexをbodyへ保存する必要がありません。",
      concepts: [
        { h: "整数が正解", p: "bodyとfoodの比較はCell同士。", code: "body []cell" },
        { h: "前回位置を一時保存", p: "移動直前の各節をvisualFromへ写す。", code: "visualFrom []cell" },
        { h: "形は隣接から導出", p: "前・現在・次で直線／角を決める。", code: "PoseForSegment" },
      ],
      lab: { title: "tick速度と補間時間を対応させる", p: "ルールは8Fごとに一マス、表示Tweenも8Fで完了させます。補間途中でも自己衝突は常に整数bodyだけで判定します。" },
      play: { title: "曲がり角と頭の滑走を見る", p: "頭はセル間を滑り、胴体は太線と丸い節で連続します。餌は脈動し、食べた瞬間だけ成長のリングと粒が出ます。" },
      codeHead: { title: "old/new二枚の盤面からDraw座標を作る", p: "Updateの離散性を捨てず、Draw座標だけLerpします。これが倉庫番やタクティクスにも続く基本です。" },
      test: { h: "角の組合せを表で試す", p: "水平、垂直、左上角などの三セルを渡し、SegmentKindとRotationを比較します。補間の0と1も厳密に確認します。" },
      failure: { h: "補間座標でセル判定しない", p: "0.4マス移動中の頭を丸めると、FPSや丸め方で自己衝突の結果が変わります。" },
      payoff: { h: "盤面ゲーム全般の標準形になる", p: "チェス駒、SRPG、ローグライクも、離散状態＋表示補間ならルールを単純に保てます。" },
      architecture: {
        name: "Discrete State, Continuous Presentation",
        problem: "滑らかさのため盤面座標をfloatにすると、同一セル・自己衝突・リプレイの答えが曖昧になります。",
        flow: ["old body cells", "commit new cells", "step Tween", "Lerp positions", "derive corners"],
        stateTitle: "visualFromは次tickまでの一時コピー",
        state: "保存対象はbodyだけ。visualFromとstepTweenは見せ終えたら次の移動で上書きできる派生状態です。",
        testTitle: "グリッド規則と形状関数を分離",
        tests: ["Progress 0/1で旧/新中心に一致", "水平・垂直はStraight", "直交する隣接はCorner"],
      },
    },
    en: {
      nav: "Discrete interpolation",
      title: "Connect a Grid Snake Smoothly",
      lead: "Snake truth is an integer cell list. Turning it into floats breaks self-collision and food tests. Copy body into visualFrom before a move, then interpolate each segment until the next tick. Derive straight or corner shape from previous/current/next cells and draw a seamless body with thick lines and circles.",
      kicker: "DISCRETE INTERPOLATION",
      question: "How can an integer board<br>look continuous?",
      explain: "Rules still jump exactly one cell. Only presentation moves old→new over several frames. SegmentPose is a pure neighbor function returning Straight / Corner, so body never stores spriteIndex.",
      concepts: [
        { h: "Integers are truth", p: "Body and food compare Cell values.", code: "body []cell" },
        { h: "Keep one old pose", p: "Copy segments to visualFrom before commit.", code: "visualFrom []cell" },
        { h: "Derive shape from neighbors", p: "Previous/current/next choose line or corner.", code: "PoseForSegment" },
      ],
      lab: { title: "Match tween duration to rule tick", p: "Rules move one cell every eight frames; presentation completes in eight. Self-collision always reads integer body during the tween." },
      play: { title: "Watch corners and the head glide", p: "The head slides between cells and the body stays connected. Food pulses; eating alone spawns growth ring and particles." },
      codeHead: { title: "Build Draw coordinates from old/new boards", p: "Keep discrete Update and Lerp only the Draw position—the same foundation used by Sokoban and tactics." },
      test: { h: "Table-test corner combinations", p: "Pass horizontal, vertical, and corner triples and compare SegmentKind and Rotation. Also test interpolation at exactly 0 and 1." },
      failure: { h: "Never test cells from interpolated coordinates", p: "Rounding a head at 0.4 cells makes collision depend on frame rate and rounding." },
      payoff: { h: "This becomes the grid-game default", p: "Chess, tactics, and roguelikes stay simple with discrete truth plus presentation interpolation." },
      architecture: {
        name: "Discrete State, Continuous Presentation",
        problem: "Float board positions make same-cell, self-collision, and replay answers ambiguous.",
        flow: ["old body cells", "commit new cells", "step Tween", "Lerp positions", "derive corners"],
        stateTitle: "visualFrom is a one-tick copy",
        state: "Only body is saved. visualFrom and stepTween are derived presentation state overwritten at the next move.",
        testTitle: "Separate grid rules from shape derivation",
        tests: ["Progress 0/1 equals old/new centers", "horizontal and vertical are Straight", "orthogonal neighbors form Corner"],
      },
    },
    code: `oldBody := append([]cell(nil), g.body...)
g.body = commitGridStep(g.body, nextHead) // integer truth
g.visualFrom = originsFor(g.body, oldBody)
g.stepTween = vfxmotion.NewTween(ruleTickFrames)

drawX := vfxmotion.Lerp(oldX, newX, g.stepTween.Progress())`,
  }),

  lesson({
    slug: "vfx-fx-shooter",
    step: "A08",
    core: "space-shooter",
    concept: { ja: "型付きイベント列と演出予算", en: "Typed event queue and FX budget" },
    hubDesc: {
      ja: "Shot・Destroyed・Damagedを列へ積み、演出dispatcherが上限内で火花や反動を選ぶ。",
      en: "Queue Shot, Destroyed, and Damaged facts; let a dispatcher spend a bounded FX budget.",
    },
    ja: {
      nav: "イベント列と予算",
      title: "大量の発射と爆発を予算内でさばく",
      lead: "シューティングでは一tickに発射・複数撃破・被弾が重なります。ルール更新はShotFired / EnemyDestroyed / PlayerDamagedを型付きQueueへ順番に積み、dispatchFXが一度だけDrainします。粒子数はBudget.Takeで上限を持ち、弾の光跡、機体の反動、爆発の波を端末性能に依存せず管理します。",
      kicker: "EVENT QUEUE + BUDGET",
      question: "同時多発する演出を、<br>どう順序と上限つきで処理する？",
      explain: "ゲーム側は「どのエフェクト関数を呼ぶか」を知りません。Kind・位置・Strengthだけを発行します。dispatcherが見せ方を選ぶため、低品質モードではBudgetを下げるだけで弾数や得点を変えずに負荷を落とせます。",
      concepts: [
        { h: "型付きイベント", p: "文字列でなくEventKindを積む。", code: "events.Push(Event{Kind: ...})" },
        { h: "一度だけDrain", p: "発行順を保って空にする。", code: "for _, e := range Drain()" },
        { h: "任意処理に上限", p: "粒子個数だけBudget.Takeで削る。", code: "n := budget.Take(wanted)" },
      ],
      lab: { title: "同tickに三種類のイベントを積む", p: "Queueの順序を保ったままdispatcherが処理し、Budgetを越えた粒だけが減ることを追います。得点イベント自体は捨てません。" },
      play: { title: "連射しながら複数の敵を壊す", p: "弾は短い光線を引き、機体は発射ごとに潰れて戻ります。撃破は橙の二重波、被弾は赤い画面フラッシュになります。" },
      codeHead: { title: "ルールはEventを発行するだけ", p: "発射処理からBurstを外し、演出の対応と負荷制御をdispatcherへ集約します。" },
      test: { h: "Queueの順序・一回性・Budgetを試す", p: "二イベントをPushしてDrain後に空になること、Limit 5へ4+4を要求すると4+1になることを確認します。" },
      failure: { h: "イベントそのものを負荷対策で捨てない", p: "得点や音まで同じBudgetで省くと、端末性能でゲーム結果が変わります。削るのは任意の粒だけです。" },
      payoff: { h: "録画・リプレイ・音声へ分配できる", p: "型付き事実の列はFX以外のconsumerも読め、デバッグログにもそのまま使えます。" },
      architecture: {
        name: "Typed Queue with an FX Budget",
        problem: "衝突ループから直接爆発を出すと、同時撃破時の順序・重複・最大負荷を制御できません。",
        flow: ["gameplay loops", "typed Queue", "Drain once", "Budget.Take", "recoil + FX"],
        stateTitle: "Queueはtickの事実、Budgetは任意描画の上限",
        state: "Queueは処理後に必ず空にします。Budgetは毎tickResetし、粒子要求だけを切り詰めます。",
        testTitle: "配達保証と負荷制御を試す",
        tests: ["Push順でDrainされる", "Drainは同じEventを二度返さない", "BudgetはLimitを越えて許可しない"],
      },
    },
    en: {
      nav: "Event queue & budget",
      title: "Dispatch Heavy Combat FX within a Budget",
      lead: "A shooter can fire, destroy several enemies, and take damage in one tick. Rules append typed ShotFired / EnemyDestroyed / PlayerDamaged events in order; dispatchFX drains once. Budget.Take caps optional particles while tracers, recoil, and shockwaves remain consistent across devices.",
      kicker: "EVENT QUEUE + BUDGET",
      question: "How do concurrent effects keep<br>order and a hard ceiling?",
      explain: "Gameplay never chooses an effect function. It emits Kind, position, and Strength. The dispatcher chooses the presentation, so a low-quality mode can lower Budget without changing bullets or score.",
      concepts: [
        { h: "Typed event", p: "Queue EventKind, not magic strings.", code: "events.Push(Event{Kind: ...})" },
        { h: "Drain once", p: "Preserve issue order, then empty.", code: "for _, e := range Drain()" },
        { h: "Cap optional work", p: "Only particle count uses Budget.Take.", code: "n := budget.Take(wanted)" },
      ],
      lab: { title: "Queue three event kinds in one tick", p: "Trace ordered dispatch and see only excess particles reduced by Budget. The gameplay facts themselves are never dropped." },
      play: { title: "Fire while destroying several enemies", p: "Bullets draw short tracers; the ship squashes on fire. Kills get orange double waves and damage gets a red screen flash." },
      codeHead: { title: "Rules only publish Events", p: "Remove Burst from firing code; centralize mapping and load control in the dispatcher." },
      test: { h: "Test order, once-only drain, and budget", p: "Push two events, verify Drain empties the queue, and verify a limit-five budget grants requests 4+4 as 4+1." },
      failure: { h: "Never budget away the event itself", p: "Dropping score or audio with particles makes game results depend on device speed. Only optional visuals are capped." },
      payoff: { h: "Fan out to recording, replay, and audio", p: "A typed fact stream supports more consumers and doubles as a debug log." },
      architecture: {
        name: "Typed Queue with an FX Budget",
        problem: "Direct explosions inside collision loops cannot control order, duplication, or worst-case load during multi-kills.",
        flow: ["gameplay loops", "typed Queue", "Drain once", "Budget.Take", "recoil + FX"],
        stateTitle: "Queue is tick facts; Budget caps optional drawing",
        state: "Queue must be empty after dispatch. Budget resets per tick and trims only particle requests.",
        testTitle: "Verify delivery and load control",
        tests: ["Drain preserves Push order", "Drain never returns an Event twice", "Budget never grants beyond Limit"],
      },
    },
    code: `g.events.Push(vfxmotion.Event{
    Kind: vfxmotion.EventEnemyDestroyed, X: e.x, Y: e.y, Strength: fx,
})

for _, event := range g.events.Drain() {
    n := g.budget.Take(wantedParticles(event))
    g.dispatch(event, n)
}`,
  }),

  lesson({
    slug: "vfx-fx-sokoban",
    step: "A09",
    core: "sokoban",
    concept: { ja: "移動を計画・確定・補間に分ける", en: "Plan, commit, then tween a move" },
    hubDesc: {
      ja: "歩行・押し・失敗をMovePlanで先読みし、盤面は即確定、絵だけマス間を移動する。",
      en: "Resolve walk, push, or block as MovePlan; commit the board now and tween only the picture.",
    },
    ja: {
      nav: "Plan→Commit→Tween",
      title: "倉庫番の一手を滑らかに見せる",
      lead: "倉庫番の盤面は一手で即座に次の整数状態へ移るのが正解です。しかしDrawも即座に次セルへ飛ばすと硬く見えます。PlanSokobanMoveで歩行・押し・失敗を副作用なしに求め、許可されたPlanを一度だけ盤面へ確定し、保存したfrom/toを9フレーム補間します。補間中は次入力を受けません。",
      kicker: "PLAN → COMMIT → TWEEN",
      question: "一手の正解を保ったまま、<br>どうマス間を移動させる？",
      explain: "MovePlanにはPlayerFrom/Toと、押す場合だけBoxIndex・BoxFrom/Toがあります。Commit後のplayer/boxesが唯一の盤面真実です。moveAnimationは同じPlanとTweenを表示用に保持し、完了時に土煙・ゴール波・クリア紙吹雪を出します。",
      concepts: [
        { h: "Planは純粋", p: "壁・箱からAllowedと移動先を返す。", code: "PlanSokobanMove(...)" },
        { h: "Commitは一回", p: "playerと必要なboxを整数セルへ更新。", code: "apply(plan)" },
        { h: "Tweenは表示だけ", p: "from/toをEaseしてDraw座標にする。", code: "Lerp(from, to, t)" },
      ],
      lab: { title: "歩行・押し・押せないをPlanで比べる", p: "壁、空き、箱の先が空き、箱の先も塞がるケースを表で試します。Allowed=falseなら盤面もTweenも作りません。" },
      play: { title: "箱を押し、足と箱の移動を追う", p: "盤面は入力時に確定しますが、海老と箱は同じProgressで滑らかに進みます。押し終わりに土煙、ゴール到達に二重波、最後の手の補間後にクリア演出が出ます。" },
      codeHead: { title: "アニメーション完了時に盤面を変更しない", p: "正解は入力時にCommit済みです。完了コールバックは演出だけを行い、二重移動を防ぎます。" },
      test: { h: "MovePlanを盤面表で網羅する", p: "空きへ歩く、壁、箱を押す、箱が箱に阻まれる、箱が壁に阻まれるを画像なしで比較します。" },
      failure: { h: "小数座標を盤面の真実にしない", p: "0.5マスの箱がどのゴールにいるか曖昧になり、Undo・保存・クリア判定が壊れます。" },
      payoff: { h: "チェス・SRPG・ターン制移動の基礎", p: "Planを作って確定し、表示だけ補間する形は、グリッド上の行動を安全に演出する標準です。" },
      architecture: {
        name: "Plan, Commit, Tween",
        problem: "盤面を毎フレーム少しずつ動かすと、補間途中の壁・箱・ゴール判定を定義しなければなりません。",
        flow: ["input direction", "pure MovePlan", "commit grid once", "input lock + Tween", "dust / goal FX"],
        stateTitle: "moveAnimationは確定済みPlanを読む",
        state: "盤面を二枚持たず、Commit後の真実とfrom/toだけを保持します。補間中の入力ロックもmove != nilで明確です。",
        testTitle: "PlanとTweenをGPUなしで保証",
        tests: ["壁へのPlanはAllowed=false", "押せる箱はBoxToを一セル先へ返す", "Tween完了前後で盤面Commitは一回だけ"],
      },
    },
    en: {
      nav: "Plan→Commit→Tween",
      title: "Show One Sokoban Move Smoothly",
      lead: "Sokoban truth should jump to the next integer board in one move. Drawing that jump immediately feels stiff. PlanSokobanMove resolves walk, push, or block without side effects; an allowed plan commits once, while its captured from/to points tween for nine frames. Input is locked during the tween.",
      kicker: "PLAN → COMMIT → TWEEN",
      question: "How can one correct move<br>travel between cells?",
      explain: "MovePlan contains PlayerFrom/To and, for a push, BoxIndex and BoxFrom/To. Committed player/boxes are the only board truth. moveAnimation retains that Plan and a Tween for drawing, then emits dust, goal wave, and clear confetti at completion.",
      concepts: [
        { h: "Plan is pure", p: "Return Allowed and destinations from walls/boxes.", code: "PlanSokobanMove(...)" },
        { h: "Commit once", p: "Update integer player and optional box.", code: "apply(plan)" },
        { h: "Tween only presentation", p: "Ease from/to into Draw positions.", code: "Lerp(from, to, t)" },
      ],
      lab: { title: "Compare walk, push, and blocked plans", p: "Table-test wall, empty, push into empty, and push into blocked. Allowed=false creates neither board mutation nor Tween." },
      play: { title: "Push a box and follow both motions", p: "The board commits on input, while Ebi and box share one smooth Progress. Dust lands at the end, goals get a double wave, and clear waits for the final tween." },
      codeHead: { title: "Do not change the board on animation completion", p: "Truth was committed at input. Completion callbacks emit presentation only, preventing a double move." },
      test: { h: "Cover MovePlan with board tables", p: "Compare empty walk, wall, valid push, box behind box, and box against wall without images." },
      failure: { h: "Never make float position board truth", p: "A half-cell box makes goal, undo, save, and clear semantics ambiguous." },
      payoff: { h: "Foundation for chess, tactics, and turns", p: "Plan, commit, then tween is the safe default for animated grid actions." },
      architecture: {
        name: "Plan, Commit, Tween",
        problem: "Moving the board a little each frame forces collision and goal rules for every in-between fraction.",
        flow: ["input direction", "pure MovePlan", "commit grid once", "input lock + Tween", "dust / goal FX"],
        stateTitle: "moveAnimation reads an already committed Plan",
        state: "Keep no second board—only committed truth and captured endpoints. move != nil is also the explicit input lock.",
        testTitle: "Guarantee Plan and Tween without a GPU",
        tests: ["a wall plan is not Allowed", "a valid push returns BoxTo one cell ahead", "board commit happens once across the Tween"],
      },
    },
    code: `plan := vfxmotion.PlanSokobanMove(player, direction, walls, boxes)
if !plan.Allowed { return }

commit(plan) // integer board truth changes once
g.move = &moveAnimation{
    plan: plan, tween: vfxmotion.NewTween(9),
}
// Draw lerps plan.PlayerFrom → plan.PlayerTo.`,
  }),

  lesson({
    slug: "vfx-fx-platform",
    step: "A10",
    core: "platformer",
    concept: { ja: "状態のエッジで姿勢機械を進める", en: "Drive poses from state edges" },
    hubDesc: {
      ja: "接地false→trueを一回のLandingにし、Idle/Run/Rise/Fall/Landを物理から選ぶ。",
      en: "Turn grounded false→true into one Landing and derive Idle/Run/Rise/Fall/Land from physics.",
    },
    ja: {
      nav: "エッジ駆動の姿勢機械",
      title: "走る・跳ぶ・着地する身体を作る",
      lead: "接地中ずっと土煙を出すのではなく、前フレームと今フレームを比較してLandedを一度だけ作ります。水平・垂直速度と接地からIdle / Run / Rise / Fallを導出し、Landingだけ短いPoseとして保持します。走りの上下動、上昇の伸び、落下の広がり、着地の潰れを当たり判定と分離します。",
      kicker: "EDGE-DRIVEN POSE MACHINE",
      question: "継続する状態から、<br>一回の着地をどう取り出す？",
      explain: "groundedは物理の状態、landedはその変化です。DetectGroundEdges(was, now)がLanded / TookOffを返し、LandedだけがShockwave・Dust・landPulseを開始します。PoseForPlatformは速度と接地を読み、表示enumを返します。",
      concepts: [
        { h: "状態とエッジ", p: "groundedのfalse→trueだけをLandedにする。", code: "DetectGroundEdges" },
        { h: "姿勢enum", p: "速度からIdle/Run/Rise/Fallを導出。", code: "PoseForPlatform" },
        { h: "短い着地Pose", p: "landPulse中だけ横へ潰す。", code: "LocomotionLand" },
      ],
      lab: { title: "接地列からイベント回数を数える", p: "false,false,true,true,falseという列でLandedとTookOffが各一回だけになることを確認します。" },
      play: { title: "走り、跳び、着地の形を見比べる", p: "走行中は小さく呼吸し、上昇で縦長、落下で横広、着地で大きく潰れます。足元の二重波と左右へ広がる土煙は着地瞬間だけです。" },
      codeHead: { title: "物理更新後にエッジを求める", p: "衝突解決後のgroundedを確定してから前状態と比較し、その後に表示Poseを選びます。" },
      test: { h: "状態列と速度表を別々に試す", p: "エッジ検出はbool列、Pose選択はvx/vy/groundedの表として検証します。" },
      failure: { h: "Draw中にwasGroundedを更新しない", p: "描画回数や画面非表示でイベント回数が変わり、テストも再現もできなくなります。" },
      payoff: { h: "入力・AI・ネット同期と独立する", p: "姿勢が物理結果から決まるため、人間操作でも敵AIでも同じ着地表現を共有できます。" },
      architecture: {
        name: "Edge-driven Pose Machine",
        problem: "grounded中に毎tick演出を呼ぶと土煙が連発し、手書きPoseフラグは物理と食い違います。",
        flow: ["resolve physics", "grounded edge", "Locomotion", "pose Tween", "squash + dust"],
        stateTitle: "永続状態は物理、短期状態はlandPulse",
        state: "Idle/Run/Rise/Fallは毎frame導出し、継続時間が必要なLandだけTweenを所有します。",
        testTitle: "一回性と姿勢対応を保証",
        tests: ["false→trueでLanded一回", "空中vy<0はRise", "landPulse中はLandが優先"],
      },
    },
    en: {
      nav: "Edge-driven pose machine",
      title: "Build a Body that Runs, Jumps, and Lands",
      lead: "Do not emit dust every grounded frame. Compare previous and current physics to create one Landed edge. Derive Idle / Run / Rise / Fall from velocity and grounded, retaining only Landing as a short pose. Run bob, rise stretch, fall width, and landing squash stay outside collision.",
      kicker: "EDGE-DRIVEN POSE MACHINE",
      question: "How does a continuous state<br>yield one landing?",
      explain: "grounded is physics state; landed is its change. DetectGroundEdges(was, now) returns Landed / TookOff. Only Landed starts Shockwave, Dust, and landPulse. PoseForPlatform reads velocity and grounded to return a display enum.",
      concepts: [
        { h: "State vs edge", p: "Only grounded false→true is Landed.", code: "DetectGroundEdges" },
        { h: "Pose enum", p: "Derive Idle/Run/Rise/Fall from physics.", code: "PoseForPlatform" },
        { h: "Short landing pose", p: "Squash only during landPulse.", code: "LocomotionLand" },
      ],
      lab: { title: "Count events from a grounded sequence", p: "For false,false,true,true,false, verify exactly one Landed and one TookOff." },
      play: { title: "Compare run, rise, fall, and land shapes", p: "Run breathes, rise stretches, fall widens, and land squashes. Double wave and bilateral dust occur only on the landing edge." },
      codeHead: { title: "Find edges after physics resolution", p: "Commit grounded after collision, compare it to the previous state, then choose presentation Pose." },
      test: { h: "Test state sequences and velocity tables", p: "Test edge detection as bool sequences and Pose mapping as vx/vy/grounded rows." },
      failure: { h: "Never update wasGrounded in Draw", p: "Draw count and hidden windows would change event count, destroying replay and tests." },
      payoff: { h: "Independent of input, AI, and networking", p: "Physics-derived poses work identically for a player, enemy AI, or synchronized remote unit." },
      architecture: {
        name: "Edge-driven Pose Machine",
        problem: "Emitting FX every grounded tick spams dust, while hand-written pose flags drift from physics.",
        flow: ["resolve physics", "grounded edge", "Locomotion", "pose Tween", "squash + dust"],
        stateTitle: "Physics owns durable state; landPulse owns the transient",
        state: "Idle/Run/Rise/Fall derive every frame. Only Land needs a Tween because it has a visible duration.",
        testTitle: "Guarantee once-only edges and pose mapping",
        tests: ["false→true emits one Landed", "airborne vy<0 selects Rise", "active landPulse overrides with Land"],
      },
    },
    code: `wasGrounded := g.onGround
g.resolvePhysics()
edges := vfxmotion.DetectGroundEdges(wasGrounded, g.onGround)
if edges.Landed {
    g.landPulse = vfxmotion.NewTween(10)
    g.fx.Dust(feetX, feetY, g.vx, 14, dust)
}
g.pose = vfxmotion.PoseForPlatform(g.vx, g.vy, g.onGround, landFrames)`,
  }),

  lesson({
    slug: "vfx-fx-dungeon",
    step: "A11",
    core: "dungeon",
    concept: { ja: "ダメージ後のリアクション時間軸", en: "A post-damage reaction timeline" },
    hubDesc: {
      ja: "HPは即減らし、HitStop→Flash→Recoverを表示専用タイムラインで再生する。",
      en: "Apply HP immediately, then play HitStop→Flash→Recover on a visual-only timeline.",
    },
    ja: {
      nav: "段階的ヒットリアクション",
      title: "一撃を時間の層に分解する",
      lead: "気持ちよい攻撃は「HPを減らして粒を出す」だけではありません。命中判定は一回でHPを即減らし、ReactionがHitStop→Flash→Recoverを順に進めます。敵は左右へ振れ、斬撃は円弧を描き、白発光から元色へ戻ります。一回の斬撃が同じ敵へ毎フレーム多段ヒットしないようslashHitも分離します。",
      kicker: "REACTION TIMELINE",
      question: "一回のダメージを、<br>複数段階の手応えにするには？",
      explain: "HPはplayの真実なので命中フレームに確定します。Reactionは表示上の停止・白発光・減衰オフセットだけを持ち、Done後は0へ戻ります。ノックバック風のOffsetを衝突座標へ使わないため、壁抜けや二重ヒットを生みません。",
      concepts: [
        { h: "数字は即時", p: "命中一回につきHPを一度だけ減らす。", code: "hp--" },
        { h: "時間を段階化", p: "HitStop / Flash / Recoverを明示。", code: "Reaction.Phase()" },
        { h: "表示オフセット", p: "揺れはDraw座標にだけ足す。", code: "x + reaction.Offset(7)" },
      ],
      lab: { title: "2F停止・5F発光・10F回復を追う", p: "Reactionを1FずつAdvanceし、期待Phase列とDone後Offset=0を比較します。" },
      play: { title: "一撃ごとの停止と白発光を感じる", p: "斬撃は扇状の三本線で走り、命中瞬間に短く止まります。敵は白く光って左右に震え、死亡時は小さな紙吹雪へ変わります。" },
      codeHead: { title: "applyDamageとstartReactionを分ける", p: "数字の適用後にReactionを作ります。Reaction完了時にもう一度HPを減らしてはいけません。" },
      test: { h: "Phase列と一回のHP差分を試す", p: "HitStop×2、Flash×5、Recover×10、Doneという順序と、斬撃一回のHP差分が1であることを別々に確認します。" },
      failure: { h: "表示ノックバックを物理へ混ぜない", p: "細かな左右振動を実座標へ足すと、壁判定・敵接触・リプレイ結果まで揺れます。" },
      payoff: { h: "コンボ・ボス・無敵時間を整理できる", p: "攻撃結果とリアクション時間が別なら、次の行動可否やキャンセル規則を名前のあるPhaseとして設計できます。" },
      architecture: {
        name: "Layered Hit-reaction Timeline",
        problem: "一つのtimerへ停止・発光・揺れをifで重ねると、境界とキャンセル規則が読めません。",
        flow: ["hit once", "apply HP", "HitStop", "Flash", "Recover / Done"],
        stateTitle: "Reactionは表示Phaseと経過だけ",
        state: "HPやaliveを所有せず、総フレームと現在FrameからPhase・Offsetを純粋に導出します。",
        testTitle: "数字と時間軸を別に保証",
        tests: ["一斬撃は同じ敵へ一回だけ命中", "Phaseが指定フレーム順に進む", "Done後Offsetは必ず0"],
      },
    },
    en: {
      nav: "Layered hit reaction",
      title: "Split One Hit into Layers of Time",
      lead: "A satisfying hit is more than reducing HP and spawning particles. Commit HP once, then let Reaction advance HitStop→Flash→Recover. The enemy shakes, the slash draws an arc, and white flash fades back. A separate slashHit prevents one swing from damaging the same target every frame.",
      kicker: "REACTION TIMELINE",
      question: "How does one damage event<br>become several layers of feel?",
      explain: "HP is play truth and commits on the hit frame. Reaction owns only visual pause, white flash, and decaying offset, returning to zero when Done. Its fake knockback never changes collision position, so it cannot tunnel through walls or re-hit.",
      concepts: [
        { h: "Numbers now", p: "Reduce HP exactly once per hit.", code: "hp--" },
        { h: "Explicit phases", p: "Name HitStop / Flash / Recover.", code: "Reaction.Phase()" },
        { h: "Visual offset", p: "Add shake only to Draw position.", code: "x + reaction.Offset(7)" },
      ],
      lab: { title: "Trace 2 stop, 5 flash, 10 recover frames", p: "Advance Reaction one frame at a time and compare the Phase sequence plus Offset=0 after Done." },
      play: { title: "Feel the pause and white flash per strike", p: "Three lines sweep as a slash arc. Contact pauses briefly; enemies flash white, shake, and become small confetti on death." },
      codeHead: { title: "Separate applyDamage from startReaction", p: "Create Reaction after committing the number. Completion must never apply damage again." },
      test: { h: "Test phase sequence and one HP delta", p: "Verify HitStop×2, Flash×5, Recover×10, Done—and separately verify one swing changes HP by one." },
      failure: { h: "Do not feed visual knockback into physics", p: "Small render shakes in real coordinates alter walls, contact, and replay results." },
      payoff: { h: "Organize combos, bosses, and invulnerability", p: "Separated results and reaction time allow named action and cancel rules." },
      architecture: {
        name: "Layered Hit-reaction Timeline",
        problem: "Stacking pause, flash, and shake as conditions on one unnamed timer hides boundaries and cancellation rules.",
        flow: ["hit once", "apply HP", "HitStop", "Flash", "Recover / Done"],
        stateTitle: "Reaction owns only phase time",
        state: "It owns no HP or alive flag; Phase and Offset derive from duration fields and current Frame.",
        testTitle: "Guarantee numbers and timeline separately",
        tests: ["one swing hits one enemy once", "phases advance for exact durations", "Offset is zero after Done"],
      },
    },
    code: `enemy.hp-- // gameplay commits once
enemy.reaction = vfxmotion.NewReaction(2, 5, 10)

// Draw only:
x := enemy.x + enemy.reaction.Offset(7)
if enemy.reaction.Phase() == vfxmotion.ReactionFlash { tint = white }`,
  }),

  lesson({
    slug: "vfx-fx-bullethell",
    step: "A12",
    core: "bullet-hell",
    concept: { ja: "複合演出を決定的な台本にする", en: "Script composite FX deterministically" },
    hubDesc: {
      ja: "ボムで弾を即削除し、Freeze→Flash→Wave→Dissolve→Confettiを再生可能な台本にする。",
      en: "Delete bombed bullets now, then replay Freeze→Flash→Wave→Dissolve→Confetti from a deterministic script.",
    },
    ja: {
      nav: "決定的エフェクト台本",
      title: "画面全体のボム演出を台本化する",
      lead: "総仕上げは複数の演出が時間差で重なるボムです。ルール側は範囲内の弾を同tickで削除し、個数を得点へ変えます。その後BombScriptが0Fの停止とフラッシュ、1Fの衝撃波、4Fの弾消散、12Fの紙吹雪をCueとして発行します。同じStrengthとFrameなら同じCue列になるため、録画・リプレイ・テストが可能です。",
      kicker: "DETERMINISTIC EFFECT SCRIPT",
      question: "大演出の順序を、<br>どう再生・検証可能にする？",
      explain: "ボムの安全性は「弾が先に消える」ことです。clearedBulletsは消散を見せる表示用コピーで、当たり判定には戻りません。EffectScriptは時刻付きCueの配列で、Current→Advanceを一回ずつ呼びます。粒子システム自身にも上限があるため大量消去で暴走しません。",
      concepts: [
        { h: "ルールを先に確定", p: "範囲内bulletsを即削除し得点化。", code: "bullets = keep" },
        { h: "時刻付きCue", p: "Freeze/Flash/Wave等をAtで並べる。", code: "BombScript(strength)" },
        { h: "表示用消散コピー", p: "削除弾を短時間だけ輪郭として描く。", code: "clearedBullets" },
      ],
      lab: { title: "同じScriptを二回走らせCue列を比べる", p: "各FrameのCurrentを記録し、同じ入力で完全一致することを確認します。乱数を使う粒の配置はCueの外側に閉じ込めます。" },
      play: { title: "弾へ近づいてボムの層を見る", p: "ボムを押した瞬間に危険弾は当たり判定から消え、停止、白橙フラッシュ、二重波、弾の輪郭消散、紙吹雪が順番に走ります。" },
      codeHead: { title: "巨大なif timerをCue表へ変える", p: "何Fに何を始めるかをデータとして並べ、dispatcherで各Cueを解釈します。" },
      test: { h: "決定性・即時削除・上限を試す", p: "同ScriptのCue列一致、doBomb直後のbullets減少、粒子MaxParts超過時のcapを別々に確認します。" },
      failure: { h: "消散中の弾を当たり判定へ戻さない", p: "半透明だから安全に見える弾で被弾すると、演出とルールが逆の情報を伝えます。" },
      payoff: { h: "必殺技・ステージ開始・ボス撃破へ", p: "複合演出をCue台本にすると、スキップ・倍速・リプレイ・演出品質差し替えを同じ事実から実装できます。" },
      architecture: {
        name: "Deterministic Composite-effect Script",
        problem: "巨大なtimer switchへ演出を足し続けると、順序・同時発火・スキップ・再生テストが崩れます。",
        flow: ["resolve bomb", "remove bullets", "BombScript", "timed Cues", "layered FX"],
        stateTitle: "ScriptはFrameと時刻付きCueだけ",
        state: "危険弾は既にbulletsから消えています。clearedBulletsは短命な表示コピーで、Script完了時にまとめて捨てます。",
        testTitle: "同じ入力なら同じ演出命令",
        tests: ["Frame 0はFreeze+Flash", "Frame 1/4/12のCue順が固定", "削除弾コピーは衝突配列へ戻らない"],
      },
    },
    en: {
      nav: "Deterministic FX script",
      title: "Script a Full-screen Bomb Sequence",
      lead: "The finale layers several effects over time. Rules delete in-range bullets in the same tick and convert their count to score. BombScript then issues freeze and flash at frame 0, shockwave at 1, bullet dissolve at 4, and confetti at 12. Equal Strength and Frame produce equal cue streams for replay and tests.",
      kicker: "DETERMINISTIC EFFECT SCRIPT",
      question: "How can a large sequence<br>be replayed and verified?",
      explain: "Bomb safety begins with bullets disappearing first. clearedBullets is a presentation copy for dissolve and never re-enters collision. EffectScript is an array of timed Cues advanced once per frame. The particle system also has a cap, preventing mass clears from exploding load.",
      concepts: [
        { h: "Commit rules first", p: "Delete in-range bullets and score now.", code: "bullets = keep" },
        { h: "Timed Cues", p: "Place Freeze/Flash/Wave at explicit frames.", code: "BombScript(strength)" },
        { h: "Dissolve copy", p: "Draw deleted bullet outlines briefly.", code: "clearedBullets" },
      ],
      lab: { title: "Run the same Script twice and compare cues", p: "Record Current at each tick and require exact equality. Any randomized particle placement remains outside the Cue contract." },
      play: { title: "Approach bullets and read the bomb layers", p: "Dangerous bullets leave collision instantly, followed by pause, white-orange flash, double wave, outline dissolve, and confetti." },
      codeHead: { title: "Replace a giant timer switch with a Cue table", p: "Store what starts on which frame as data and interpret each Cue in a dispatcher." },
      test: { h: "Test determinism, immediate deletion, and cap", p: "Separately verify equal Cue streams, fewer bullets immediately after doBomb, and MaxParts enforcement." },
      failure: { h: "Never return dissolving bullets to collision", p: "A translucent “safe-looking” bullet that still hurts communicates the opposite of the rule." },
      payoff: { h: "Use for supers, intros, and boss deaths", p: "Cue scripts make skip, fast-forward, replay, and quality-specific presentation possible from the same facts." },
      architecture: {
        name: "Deterministic Composite-effect Script",
        problem: "Growing a giant timer switch destroys ordering, simultaneous cues, skipping, and replay tests.",
        flow: ["resolve bomb", "remove bullets", "BombScript", "timed Cues", "layered FX"],
        stateTitle: "Script owns only Frame and timed Cues",
        state: "Dangerous bullets are already absent. clearedBullets is a short-lived display copy discarded when Script completes.",
        testTitle: "Same input yields the same visual commands",
        tests: ["Frame 0 emits Freeze+Flash", "Frames 1/4/12 keep fixed cue order", "deleted bullet copies never return to collision"],
      },
    },
    code: `g.bullets, g.clearedBullets = resolveBomb(g.bullets, center, radius)
g.bombScript = vfxmotion.BombScript(strength)

for _, cue := range g.bombScript.Current() {
    g.dispatchCue(cue) // freeze, flash, wave, dissolve, confetti
}
g.bombScript.Advance()`,
  }),
];
