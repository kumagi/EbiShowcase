export const presentationPatterns = {
  A01: { slug: "vfx-fx-tap", ja: "不変の結果イベント", en: "immutable result event" },
  A02: { slug: "vfx-fx-meter", ja: "型付き演出レシピ", en: "typed presentation recipe" },
  A03: { slug: "vfx-fx-catch", ja: "表示代理の寿命", en: "visual-proxy lifecycle" },
  A04: { slug: "vfx-fx-flight", ja: "物理から導く姿勢", en: "pose derived from physics" },
  A05: { slug: "vfx-fx-pong", ja: "上限つき描画履歴", en: "bounded presentation history" },
  A06: { slug: "vfx-fx-breakout", ja: "削除前スナップショット", en: "pre-delete snapshot" },
  A07: { slug: "vfx-fx-snake", ja: "離散状態の補間", en: "discrete-state interpolation" },
  A08: { slug: "vfx-fx-shooter", ja: "イベント列と演出予算", en: "event queue and FX budget" },
  A09: { slug: "vfx-fx-sokoban", ja: "Plan→Commit→Tween", en: "Plan→Commit→Tween" },
  A10: { slug: "vfx-fx-platform", ja: "エッジ駆動の姿勢機械", en: "edge-driven pose machine" },
  A11: { slug: "vfx-fx-dungeon", ja: "ヒットリアクション時間軸", en: "hit-reaction timeline" },
  A12: { slug: "vfx-fx-bullethell", ja: "決定的エフェクト台本", en: "deterministic effect script" },
};

const track = (id, patterns, ja, en) => ({ id, patterns, ja, en });

// Each specialization assumes the selected A01–A12 tools. "focus" is the new
// genre rule; "motion" is the presentation baseline that must already be
// visible in demos; "test" names the seam that keeps the combination safe.
export const genrePresentationMap = [
  track("platformer", ["A10", "A04", "A06"], {
    name: "横スクロールアクション",
    focus: "ここで新しく学ぶのは、連続衝突、足場、カメラ、ステージデータです。Update / Drawや着地演出を最初から教え直しません。",
    motion: "Idle / Run / Rise / Fall / Land、敵の予備動作・接触・回復、コインや敵が消えた後の余韻までを最初の完成基準にします。",
    test: "物理と接地エッジを純粋ロジックで試し、Poseと粒子を外しても同じ座標・衝突結果になること。",
  }, {
    name: "Side-scrolling platformer",
    focus: "New work is continuous collision, platforms, camera, and stage data—not relearning Update / Draw or landing feedback.",
    motion: "Idle / Run / Rise / Fall / Land, enemy anticipation/contact/recovery, and pickup/death afterlife are the baseline.",
    test: "Test physics and grounded edges as pure logic; removing poses and particles must preserve positions and collisions.",
  }),
  track("survivors", ["A08", "A11", "A06"], {
    name: "大量敵サバイバル",
    focus: "新しい課題は大量エンティティ、空間探索、自動攻撃、経験値と強化選択です。",
    motion: "攻撃イベントを列で配り、敵死亡は即削除＋表示スナップショット、被弾は段階的リアクションとして扱います。",
    test: "敵数が増えてもイベント一回性と粒子上限が守られ、演出品質で撃破数が変わらないこと。",
  }, {
    name: "Arena survivors",
    focus: "New work is entity scale, spatial queries, auto attacks, XP, and upgrade choices.",
    motion: "Dispatch attacks through events, show death from snapshots, and stage damage as a reaction timeline.",
    test: "Event once-only delivery and FX caps survive large crowds; quality settings never change kill count.",
  }),
  track("clicker", ["A02", "A03", "A12"], {
    name: "放置クリッカー",
    focus: "新しい課題は指数的価格、生産ライン、保存、離席時間の計算です。",
    motion: "購入結果を演出レシピへ写し、生産物は表示代理としてHUDへ運び、大きな節目はCue台本で祝います。",
    test: "経済計算と経過時間を純粋に試し、数字の増加が紙吹雪の長さや描画FPSに依存しないこと。",
  }, {
    name: "Idle clicker",
    focus: "New work is exponential prices, production lines, persistence, and offline elapsed time.",
    motion: "Map purchases to recipes, fly production as visual proxies, and celebrate milestones with cue scripts.",
    test: "Economy and elapsed-time math stay pure; production never depends on confetti duration or render FPS.",
  }),
  track("rpg", ["A11", "A02", "A12"], {
    name: "コマンドRPG",
    focus: "新しい課題はターン、コマンド、状態効果、会話フラグ、クエスト進行です。",
    motion: "行動を予備→接触→回復の時間軸にし、技ごとの表示レシピと勝利・レベルアップの台本を組みます。",
    test: "ダメージ式とターン遷移を演出時間から切り離し、行動一回がHPへ一度だけ反映されること。",
  }, {
    name: "Command RPG",
    focus: "New work is turns, commands, statuses, dialogue flags, and quest progression.",
    motion: "Stage actions as anticipation→contact→recovery, use skill recipes, and script victory/level-up sequences.",
    test: "Damage and turn transitions stay independent of presentation time; one action applies HP once.",
  }),
  track("fighting", ["A11", "A04", "A05"], {
    name: "2D格闘ゲーム",
    focus: "新しい課題は押し合い・喰らい・攻撃判定、フレームデータ、キャンセル、間合いAIです。",
    motion: "立ち・歩き・始動・持続・硬直・喰らいを名前のあるPose/Phaseにし、拳や足の軌跡は上限つき履歴で描きます。",
    test: "技の判定発生F・持続F・硬直Fを表で試し、見た目の補間が当たり判定の時刻を変えないこと。",
  }, {
    name: "2D fighting game",
    focus: "New work is push/hurt/attack boxes, frame data, cancels, and spacing AI.",
    motion: "Name idle, walk, startup, active, recovery, and hurt poses; draw limb arcs from bounded history.",
    test: "Table-test startup/active/recovery frames; visual interpolation never moves hitbox timing.",
  }),
  track("merge-physics", ["A01", "A06", "A08"], {
    name: "合体物理パズル",
    focus: "新しい課題は円の衝突解決、めり込み修正、反発、重力、合体サイズ曲線です。",
    motion: "MergeResultを一回の事実にし、消えた二体のSnapshotから新しい一体への収束を描き、同時合体は予算つき列で処理します。",
    test: "同じ初期状態とtick数で同じ合体結果になり、粒数を0にしても半径・得点・ゲームオーバーが一致すること。",
  }, {
    name: "Merge physics puzzle",
    focus: "New work is circle resolution, penetration correction, bounce, gravity, and tier-size curves.",
    motion: "Create one MergeResult, animate two deleted snapshots into the new body, and budget simultaneous merge events.",
    test: "Equal initial state and ticks yield equal merges; zero particles preserve radii, score, and game-over.",
  }),
  track("deckbuilder", ["A03", "A11", "A12"], {
    name: "デッキ構築カードゲーム",
    focus: "新しい課題は山札・手札・捨て札、エネルギー、敵予告、報酬と分岐です。",
    motion: "使ったカードは表示代理として対象へ飛び、攻撃はリアクション時間軸、報酬・階層遷移はCue台本にします。",
    test: "カード移動は先にゾーン間で確定し、飛行中の表示代理がもう一度使用・捨て札化されないこと。",
  }, {
    name: "Deckbuilder",
    focus: "New work is draw/hand/discard zones, energy, intents, rewards, and routes.",
    motion: "Fly played cards as proxies, resolve attacks on reaction timelines, and script rewards/floor transitions.",
    test: "Commit zone movement first; a flying visual proxy can never be played or discarded twice.",
  }),
  track("match3", ["A09", "A06", "A12"], {
    name: "3マッチパズル",
    focus: "新しい課題は交換可否、並び探索、重力、補充、連鎖と特殊ピースです。",
    motion: "交換はPlan→Commit→Tween、消去ピースはSnapshot、連鎖は段階Cueとして再生します。",
    test: "盤面の各Phase後に不変条件を試し、Tween途中の小数座標でmatch探索を行わないこと。",
  }, {
    name: "Match-three puzzle",
    focus: "New work is swap validity, match search, gravity, refill, cascades, and specials.",
    motion: "Use Plan→Commit→Tween for swaps, snapshots for clears, and timed cues for cascade layers.",
    test: "Check board invariants after each phase; never search matches from in-between tween coordinates.",
  }),
  track("falling-blocks", ["A09", "A06", "A12"], {
    name: "落ちものブロックパズル",
    focus: "新しい課題は回転、壁蹴り、固定、ライン判定、ホールドと落下速度です。",
    motion: "盤面はライン消去結果へ即確定し、旧行Snapshotを点滅・消散、上の行の落下をTween、得点をCueで見せます。",
    test: "回転候補とライン圧縮を純粋関数で試し、消去アニメーション中も確定盤面が一意であること。",
  }, {
    name: "Falling-block puzzle",
    focus: "New work is rotation, kicks, lock, line detection, hold, and fall speed.",
    motion: "Commit the cleared board, flash old-row snapshots, tween the fall, and show scoring through cues.",
    test: "Pure-test rotation candidates and line compression; committed board remains unique during animation.",
  }),
  track("slingshot", ["A05", "A11", "A08"], {
    name: "ひっぱり反射アクション",
    focus: "新しい課題はドラッグ照準、円反射、減速、接触攻撃、味方連携です。",
    motion: "予測線は上限つきサンプル、衝突はヒットリアクション、同tickの多重接触はイベント列で一回ずつ処理します。",
    test: "予測軌道を消しても実軌道が同じで、同じ敵IDへの接触が一発射中に規定回数だけになること。",
  }, {
    name: "Slingshot action",
    focus: "New work is drag aiming, circle reflection, damping, contact attacks, and ally links.",
    motion: "Use bounded samples for prediction, reactions for impacts, and a queue for same-tick contacts.",
    test: "Hiding prediction preserves the real trajectory; each enemy ID receives only the intended contacts per shot.",
  }),
  track("sandbox", ["A06", "A08", "A10"], {
    name: "2Dブロックサンドボックス",
    focus: "新しい課題はタイル世界、採掘・設置、チャンク、クラフト、生物、保存です。",
    motion: "破壊タイルはSnapshotから破片へ、採掘は段階Pose、同時ドロップは予算つきイベント列で扱います。",
    test: "ワールド配列の変更とインベントリ差分を純粋に試し、破片が消えても採掘結果が残ること。",
  }, {
    name: "2D sandbox",
    focus: "New work is tile worlds, mining/placing, chunks, crafting, creatures, and saves.",
    motion: "Turn deleted tiles into snapshot debris, stage mining poses, and budget simultaneous drop events.",
    test: "Pure-test world and inventory deltas; mining results remain after every shard expires.",
  }),
  track("monster-collection", ["A02", "A03", "A11"], {
    name: "モンスター収集RPG",
    focus: "新しい課題は種族データ、技、相性、捕獲率、パーティ、成長と図鑑です。",
    motion: "技は型付きレシピ、捕獲球は表示代理、攻撃と捕獲失敗は段階リアクションにします。",
    test: "相性・捕獲率・成長をseed固定で試し、球の飛行時間が捕獲結果へ影響しないこと。",
  }, {
    name: "Monster-collection RPG",
    focus: "New work is species data, moves, affinity, capture odds, party, growth, and dex.",
    motion: "Map moves to recipes, fly capture orbs as proxies, and stage attacks/failures as reactions.",
    test: "Seed-test affinity, capture, and growth; orb travel time never changes capture outcome.",
  }),
  track("falling-pairs", ["A07", "A06", "A12"], {
    name: "落下ペア連鎖パズル",
    focus: "新しい課題は二個組回転、連結探索、重力、再探索、連鎖得点とおじゃまです。",
    motion: "整数盤面を補間し、消去群をSnapshotへ移し、連鎖数ごとの消去→落下→再探索をCue台本にします。",
    test: "連結・重力・連鎖を盤面単位で試し、表示の消去時間が連鎖数を変えないこと。",
  }, {
    name: "Falling-pair chain puzzle",
    focus: "New work is pair rotation, connectivity, gravity, rescans, chain score, and garbage.",
    motion: "Interpolate the integer board, snapshot cleared groups, and script clear→fall→rescan cues.",
    test: "Test connectivity, gravity, and chains board-by-board; clear duration never changes chain count.",
  }),
  track("maze-chase", ["A07", "A10", "A11"], {
    name: "迷路追跡ゲーム",
    focus: "新しい課題はタイル移動、先行入力、交差点、敵モード、経路選択です。",
    motion: "整数セル間を補間し、交差点到達をエッジとして入力予約を適用し、被弾をリアクション時間軸にします。",
    test: "経路選択と予約入力を整数グリッドで試し、補間座標がAIの次セル選択へ入らないこと。",
  }, {
    name: "Maze chase",
    focus: "New work is tile movement, buffered turns, intersections, enemy modes, and path choice.",
    motion: "Interpolate cells, treat intersection arrival as an edge, and stage damage on a reaction timeline.",
    test: "Test routing and input buffers on integers; interpolated positions never choose an AI cell.",
  }),
  track("bomb-maze", ["A12", "A06", "A08"], {
    name: "爆弾迷路アクション",
    focus: "新しい課題は時限爆弾、壁で止まる十字爆風、壊れる壁、連鎖、アイテムです。",
    motion: "爆発をFuse→Flash→Flame→DebrisのCue台本にし、壁はSnapshot、連鎖爆発は型付き列で処理します。",
    test: "爆風セル集合と連鎖順を純粋に試し、画面揺れや炎の寿命で当たり判定が増減しないこと。",
  }, {
    name: "Bomb maze",
    focus: "New work is timed bombs, wall-stopped cross blasts, breakables, chains, and items.",
    motion: "Script Fuse→Flash→Flame→Debris, snapshot walls, and queue chain explosions.",
    test: "Pure-test blast cells and chain order; shake and flame lifetime never alter collision.",
  }),
  track("reversi", ["A07", "A09", "A02"], {
    name: "リバーシ",
    focus: "新しい課題は合法手、挟み探索、複数方向反転、パス、盤面評価CPUです。",
    motion: "手をPlanとして確定し、反転石を時間差Tween、CPUの考え方を表示レシピで可視化します。",
    test: "合法手と反転集合を盤面表で試し、石が薄く見える途中でも所有者は確定済みであること。",
  }, {
    name: "Reversi",
    focus: "New work is legal moves, directional capture, flips, passes, and CPU evaluation.",
    motion: "Commit a move plan, stagger flip tweens, and visualize CPU thinking through recipes.",
    test: "Board-test legal moves and flip sets; ownership is committed even while a stone looks half-flipped.",
  }),
  track("tactics", ["A09", "A11", "A08"], {
    name: "戦略SLG",
    focus: "新しい課題は移動コスト、到達範囲、射程、行動済み、敵思考とターン進行です。",
    motion: "移動はPlan→Commit→Tween、攻撃はリアクション時間軸、敵味方の行動列は型付きQueueで同じ再生器へ渡します。",
    test: "到達範囲・攻撃結果・行動済み遷移を純粋に試し、移動アニメーション完了で二度Commitしないこと。",
  }, {
    name: "Tactical RPG",
    focus: "New work is movement cost, reach, range, acted state, enemy planning, and turns.",
    motion: "Move with Plan→Commit→Tween, attack through reactions, and feed ally/enemy action queues to one player.",
    test: "Pure-test reach, combat, and acted transitions; animation completion never commits movement twice.",
  }),
  track("active-rpg", ["A08", "A11", "A12"], {
    name: "アクティブ戦闘RPG",
    focus: "新しい課題は速度ゲージ、READY列、入力待ち、敵行動、戦闘継続です。",
    motion: "READYを型付き列、各行動をリアクション時間軸、必殺技と戦闘終了をCue台本で再生します。",
    test: "ゲージ加算とREADY順を演出停止から分離し、一行動の効果が接触Cueで一度だけ適用されること。",
  }, {
    name: "Active-gauge RPG",
    focus: "New work is speed gauges, READY order, input waits, enemy actions, and encounter continuity.",
    motion: "Queue READY actors, stage actions as reactions, and script supers and battle endings.",
    test: "Separate gauge order from presentation pauses; each action applies its effect once at contact.",
  }),
  track("visual-novel", ["A12", "A04", "A02"], {
    name: "会話演出ゲーム",
    focus: "新しい課題はシーンデータ、文字送り、選択肢、条件分岐、好感度とエンディングです。",
    motion: "台詞・登場・表情・選択肢をCue台本で順序化し、立ち絵姿勢を状態から導出、選択結果を表示レシピへ写します。",
    test: "同じフラグと選択で同じ次sceneになり、文字送り速度や立ち絵Tweenが分岐結果を変えないこと。",
  }, {
    name: "Visual novel",
    focus: "New work is scene data, typewriter text, choices, conditions, affinity, and endings.",
    motion: "Order dialogue/entrance/expression/choice through cues, derive portrait pose, and map outcomes to recipes.",
    test: "Equal flags and choice yield the same next scene; type speed and portrait tween never alter branching.",
  }),
  track("racing", ["A04", "A05", "A10"], {
    name: "レーシングゲーム",
    focus: "新しい課題は加速、摩擦、旋回、路面、ゲート順、周回とライバル走行線です。",
    motion: "速度から車体姿勢を導出し、タイヤ跡を上限つき履歴、路面進入・ゲート通過をエッジイベントにします。",
    test: "車両物理とゲート順を純粋に試し、タイヤ跡の長さや車体傾きでラップ判定が変わらないこと。",
  }, {
    name: "Top-down racing",
    focus: "New work is acceleration, friction, steering, surfaces, ordered gates, laps, and rival lines.",
    motion: "Derive vehicle pose from velocity, keep skid marks in bounded history, and edge-trigger gates/surfaces.",
    test: "Pure-test vehicle physics and gate order; skid length and visual lean never alter lap results.",
  }),
  track("metroidvania", ["A10", "A11", "A06"], {
    name: "メトロイドヴァニア",
    focus: "新しい課題は部屋カメラ、探索記録、能力ゲート、戻り道、永続アンロックです。",
    motion: "移動Poseは物理エッジ、戦闘はリアクション、壊れる障害物や敵死亡はSnapshotの余韻で描きます。",
    test: "能力所持と部屋接続をデータで試し、カメラTweenや死亡演出が解放状態へ混ざらないこと。",
  }, {
    name: "Metroidvania",
    focus: "New work is room cameras, exploration memory, ability gates, backtracking, and persistent unlocks.",
    motion: "Drive movement poses from physics edges, combat from reactions, and destruction from snapshots.",
    test: "Data-test abilities and room links; camera tweens and death effects never enter unlock state.",
  }),
  track("raycaster", ["A04", "A05", "A11"], {
    name: "レイキャスト迷路FPS",
    focus: "新しい課題はDDA、壁投影、魚眼補正、深度、敵スプライトと射撃です。",
    motion: "速度から武器・歩行bobを導出し、弾道履歴を制限し、命中と被弾をリアクション時間軸にします。",
    test: "光線距離と遮蔽を純粋に試し、銃の反動や画面フラッシュで命中距離が変わらないこと。",
  }, {
    name: "Ray-cast maze FPS",
    focus: "New work is DDA, wall projection, fisheye correction, depth, enemy sprites, and shooting.",
    motion: "Derive weapon/walk bob from motion, bound shot history, and stage hit/damage reactions.",
    test: "Pure-test ray distance and occlusion; recoil and screen flash never alter hit range.",
  }),
  track("rhythm", ["A02", "A12", "A05"], {
    name: "リズムゲーム",
    focus: "新しい課題は音楽時刻、判定窓、譜面、長押し、連打、入力補正です。",
    motion: "判定を型付きレシピへ写し、曲と結果画面をCue台本、最近の入力差を上限つき履歴として可視化します。",
    test: "音源なしの仮想frameで判定境界を試し、粒子や画面揺れが譜面時計を停止・加速しないこと。",
  }, {
    name: "Rhythm game",
    focus: "New work is music time, judgement windows, charts, holds, rolls, and input offset.",
    motion: "Map grades to recipes, script songs/results, and visualize recent timing deltas as bounded history.",
    test: "Test judgement boundaries on a virtual frame clock; particles and shake never pause or advance chart time.",
  }),
  track("tower-defense", ["A08", "A03", "A11"], {
    name: "タワーディフェンス",
    focus: "新しい課題は経路、射程、標的優先、配置、資金、波、強化とボスです。",
    motion: "射撃・命中・撃破を予算つきイベント列、弾を表示代理、敵被弾をリアクションとして扱います。",
    test: "標的選択とダメージを純粋に試し、弾の飛行時間や粒子上限で実際の命中tickが変わらないこと。",
  }, {
    name: "Tower defense",
    focus: "New work is paths, range, targeting, placement, economy, waves, upgrades, and bosses.",
    motion: "Queue shots/hits/kills within a budget, render projectiles as proxies, and stage enemy reactions.",
    test: "Pure-test targeting and damage; projectile travel and FX caps never alter the committed hit tick.",
  }),
  track("topdown-adventure", ["A10", "A11", "A12"], {
    name: "見下ろし型アクションアドベンチャー",
    focus: "新しい課題は8方向移動、剣の範囲、無敵時間、部屋封鎖、鍵・道具・ボスPhaseです。",
    motion: "移動と攻撃Poseをエッジで開始し、命中をリアクション、部屋遷移・ボス大技・勝利をCue台本にします。",
    test: "剣判定・無敵時間・部屋条件を純粋に試し、斬撃線やカメラ揺れが敵HPや鍵状態を変えないこと。",
  }, {
    name: "Top-down action adventure",
    focus: "New work is eight-way movement, sword reach, i-frames, room seals, keys/tools, and boss phases.",
    motion: "Edge-trigger movement/attack poses, stage hit reactions, and script rooms, boss supers, and victory.",
    test: "Pure-test sword hits, i-frames, and room conditions; slash art and camera shake never change HP or keys.",
  }),
];

export function validateGenrePresentationMap() {
  const ids = new Set();
  for (const entry of genrePresentationMap) {
    if (ids.has(entry.id)) throw new Error(`duplicate genre presentation entry: ${entry.id}`);
    ids.add(entry.id);
    if (entry.patterns.length < 2 || entry.patterns.length > 3) {
      throw new Error(`${entry.id}: expected 2–3 foundation patterns`);
    }
    for (const pattern of entry.patterns) {
      if (!presentationPatterns[pattern]) throw new Error(`${entry.id}: unknown pattern ${pattern}`);
    }
    for (const lang of ["ja", "en"]) {
      for (const field of ["name", "focus", "motion", "test"]) {
        if (!entry[lang][field]) throw new Error(`${entry.id}: missing ${lang}.${field}`);
      }
    }
  }
  return ids;
}
