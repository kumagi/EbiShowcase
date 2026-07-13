// Simple optical tricks: large visual payoff without shaders or new art.
export const illusionLessons = [
  {
    slug: "vfx-squash", tier: "basic", step: "14", stars: "★★☆☆☆", labKind: "squash",
    concept: {ja: "GeoM.Scale + 足元の軸", en: "GeoM.Scale + foot anchor"},
    hubDesc: {ja: "1枚絵を潰して伸ばし、予備動作と着地の重さを作ります。", en: "Squash and stretch one sprite to sell anticipation and weight."},
    ja: {
      navConcept: "潰す・伸ばす", title: "スクワッシュ＆ストレッチ", lead: "新しいコマを描かなくても、跳ぶ前に横へ潰し、飛ぶ瞬間に縦へ伸ばすだけで、海老・天次郎の体が柔らかく見えます。大事なのは中心ではなく足元を動かさないことです。",
      deepEyebrow: "OPTICAL TRICK / SQUASH", deepH: "1枚絵に<br>重さを足す", deepLead: "幅を増やしたぶん高さを減らすと、体積を保ったように見えます。足元を軸にすれば地面へめり込まず、潰れ→伸び→戻る順番だけで予備動作と反動が伝わります。",
      concepts: [{h:"潰す",p:"跳ぶ前に幅を増やし、高さを減らします。",code:"Scale(1+s, 1-s)"},{h:"伸ばす",p:"動き出す瞬間は反対へ細長く。",code:"s = -s"},{h:"足元を固定",p:"中心ではなく足を軸に変形します。",code:"Translate(-w/2, -h)"}],
      lab:{eyebrow:"TRY IT / SHAPE",title:"潰れ方を切り替えよう",p:"潰れ・通常・伸びを切り替え、足元が同じ場所に残ることを見比べます。"},
      codeHead:{eyebrow:"IN THE WASM",h:"Scaleは2つの数",p:"横と縦を逆方向へ変え、最後に足元の座標へ戻します。"},
      whys:[{eyebrow:"WHY IT WORKS",h:"目は形の変化を力と読む",p:"急な伸びは速さ、潰れは重さに見えます。"},{eyebrow:"WHY ANCHOR",h:"足が滑らない",p:"接地点が固定されると、柔らかくても地面は固く見えます。"},{eyebrow:"TRY NEXT",h:"ボールにも使う",p:"壁へ当たる方向だけ潰して、反射へ弾みを足そう。"}],
    },
    en: {
      navConcept:"Squash & stretch",title:"Squash & Stretch",lead:"Without drawing another frame, squash Ebi Tenjiroh before a jump and stretch on takeoff. Anchor the feet instead of the center so the body feels soft while the ground stays solid.",
      deepEyebrow:"OPTICAL TRICK / SQUASH",deepH:"Give one sprite<br>weight",deepLead:"Increase width while reducing height to suggest preserved volume. A foot anchor prevents sinking. Squash → stretch → settle communicates anticipation and recoil.",
      concepts:[{h:"Squash",p:"Wider and shorter before motion.",code:"Scale(1+s, 1-s)"},{h:"Stretch",p:"Flip the shape on takeoff.",code:"s = -s"},{h:"Foot anchor",p:"Transform around the feet, not center.",code:"Translate(-w/2, -h)"}],
      lab:{eyebrow:"TRY IT / SHAPE",title:"Switch the body shape",p:"Compare squash, neutral, and stretch while the feet stay planted."},codeHead:{eyebrow:"IN THE WASM",h:"Scale is two numbers",p:"Drive horizontal and vertical in opposite directions, then place the feet."},
      whys:[{eyebrow:"WHY IT WORKS",h:"Shape reads as force",p:"A sharp stretch says speed; a squash says weight."},{eyebrow:"WHY ANCHOR",h:"No sliding feet",p:"A stable contact point keeps the ground convincing."},{eyebrow:"TRY NEXT",h:"Use it on a ball",p:"Squash along the collision axis before a bounce."}],
    },
    code:`s := math.Sin(phase) * amount
op.GeoM.Translate(-w/2, -h) // 足元を原点へ
op.GeoM.Scale(base*(1+s), base*(1-s))
op.GeoM.Translate(feetX, feetY)
screen.DrawImage(sprite, op)`,
  },
  {
    slug:"vfx-outline",tier:"basic",step:"15",stars:"★★☆☆☆",labKind:"outline",
    concept:{ja:"ColorScale + 多重描画",en:"ColorScale + repeated draws"},hubDesc:{ja:"黒いずらしコピーを下へ敷き、輪郭と影を作ります。",en:"Place tinted offset copies underneath to make an outline and shadow."},
    ja:{navConcept:"輪郭と影",title:"多重描画で輪郭を作る",lead:"同じ画像を黒く染め、少しずつ方向を変えて先に描きます。最後に本体を重ねるだけで、明るい背景でもキャラクターがくっきり読めます。",
      deepEyebrow:"OPTICAL TRICK / OUTLINE",deepH:"新しい画像なしで<br>縁を生やす",deepLead:"上下左右、または円周上へ黒いコピーを数pxずらして描きます。本体が中央を隠すため、外にはみ出した部分だけが輪郭として残ります。当たり判定は元のままです。",
      concepts:[{h:"染める",p:"コピーを黒または白一色にします。",code:"ScaleWithColor"},{h:"ずらす",p:"円周上の方向へ数px動かします。",code:"cos(a), sin(a)"},{h:"本体を重ねる",p:"最後に普通の画像で中央を隠します。",code:"DrawImage(main)"}],lab:{eyebrow:"TRY IT / OUTLINE",title:"太さと色を変えよう",p:"輪郭の太さと暗色・明色を切り替え、背景からの読みやすさを比べます。"},codeHead:{eyebrow:"IN THE WASM",h:"8枚ずらして1枚重ねる",p:"画像の追加や加工なし。DrawImageの順番だけで作れます。"},whys:[{eyebrow:"WHY IT WORKS",h:"本体が中央を隠す",p:"外側へ出たコピーだけが縁になります。"},{eyebrow:"WHY READABILITY",h:"背景に負けない",p:"明暗の境界ができ、スマホの小画面でも形を追えます。"},{eyebrow:"TRY NEXT",h:"影へ変える",p:"下方向だけ大きくずらし、薄くすれば落ち影になります。"}]},
    en:{navConcept:"Outline & shadow",title:"Outline with Repeated Draws",lead:"Tint the same sprite, draw offset copies first, then cover the center with the original. The character stays readable even on a bright background.",deepEyebrow:"OPTICAL TRICK / OUTLINE",deepH:"Grow an edge<br>without new art",deepLead:"Draw dark copies a few pixels around a circle. The original covers their centers, leaving only the outside as an outline. Collision remains unchanged.",concepts:[{h:"Tint",p:"Make each copy black or white.",code:"ScaleWithColor"},{h:"Offset",p:"Move copies around a circle.",code:"cos(a), sin(a)"},{h:"Cover",p:"Draw the normal sprite last.",code:"DrawImage(main)"}],lab:{eyebrow:"TRY IT / OUTLINE",title:"Change width and color",p:"Compare dark and light outlines at several widths."},codeHead:{eyebrow:"IN THE WASM",h:"Eight copies, then one",p:"No new asset—only tint, offsets, and draw order."},whys:[{eyebrow:"WHY IT WORKS",h:"The body hides the middle",p:"Only the outside of each copy remains."},{eyebrow:"WHY READABILITY",h:"Separate from backgrounds",p:"A value edge stays clear on small screens."},{eyebrow:"TRY NEXT",h:"Turn it into a shadow",p:"Offset only downward and lower alpha."}]},
    code:`for i := 0; i < copies; i++ {
    a := float64(i) * 2 * math.Pi / float64(copies)
    drawTinted(math.Cos(a)*width, math.Sin(a)*width, outline)
}
drawSprite(0, 0) // 本体は最後`,
  },
  {
    slug:"vfx-impact-lines",tier:"basic",step:"16",stars:"★★★☆☆",labKind:"impact",
    concept:{ja:"StrokeCircle + StrokeLine",en:"StrokeCircle + StrokeLine"},hubDesc:{ja:"広がる輪と集中線だけで、必殺技の衝撃を作ります。",en:"Expanding rings and focus lines turn a point into a special-move impact."},
    ja:{navConcept:"衝撃波と集中線",title:"衝撃波リングと集中線",lead:"タップした点から円を広げ、画面の外へ向かう短い線を一瞬だけ描きます。画像が1枚もなくても、中心へ力が集まった必殺技に見えます。",deepEyebrow:"OPTICAL TRICK / IMPACT",deepH:"円と線だけで<br>画面を殴る",deepLead:"輪は時間とともに半径を増やし、透明へ消します。集中線は中心近くを空け、外側だけ描くと視線が衝突点へ向かいます。寿命は短いほど鋭く見えます。",concepts:[{h:"広げる",p:"ageに速さを掛けて半径にします。",code:"radius := age*speed"},{h:"集める",p:"同じ中心から角度違いの線を描きます。",code:"angle := i*2π/n"},{h:"消す",p:"寿命の終わりへ透明にします。",code:"alpha := 1-age/life"}],lab:{eyebrow:"TRY IT / IMPACT",title:"衝撃を発生させよう",p:"中心を押して輪と集中線を再生し、線の数を切り替えます。"},codeHead:{eyebrow:"IN THE WASM",h:"同じ中心を共有する",p:"円と線が同じ座標を使うだけで、ひとつの強い衝撃にまとまります。"},whys:[{eyebrow:"WHY A RING",h:"力が広がって見える",p:"半径の増加を衝撃の伝わりとして読みます。"},{eyebrow:"WHY LINES",h:"視線が中心へ集まる",p:"放射状の線が重要な点を指します。"},{eyebrow:"TRY NEXT",h:"2フレーム止める",p:"ヒットストップを組み合わせるとさらに重くなります。"}]},
    en:{navConcept:"Shockwave & focus lines",title:"Shockwave Rings & Focus Lines",lead:"Expand circles from a tapped point and flash short radial lines. With no image assets, the center already reads like a special-move impact.",deepEyebrow:"OPTICAL TRICK / IMPACT",deepH:"Punch the screen<br>with circles and lines",deepLead:"Grow ring radius with age and fade it out. Leave a gap near the center of each focus line so the eye is pulled toward the collision. A short lifetime feels sharp.",concepts:[{h:"Expand",p:"Multiply age by speed for radius.",code:"radius := age*speed"},{h:"Focus",p:"Draw many angles from one center.",code:"angle := i*2π/n"},{h:"Fade",p:"Lose alpha toward end of life.",code:"alpha := 1-age/life"}],lab:{eyebrow:"TRY IT / IMPACT",title:"Trigger the impact",p:"Replay the ring and change the number of focus lines."},codeHead:{eyebrow:"IN THE WASM",h:"Share one center",p:"Circles and lines using one coordinate merge into one strong hit."},whys:[{eyebrow:"WHY A RING",h:"Force appears to spread",p:"Growing radius reads as transmitted impact."},{eyebrow:"WHY LINES",h:"The eye finds the center",p:"Radial lines point at the important spot."},{eyebrow:"TRY NEXT",h:"Freeze for two frames",p:"Pair it with hit stop for extra weight."}]},
    code:`radius := age * speed
vector.StrokeCircle(screen, x, y, radius, width, color, true)
for i := 0; i < rays; i++ {
    angle := float64(i) * 2 * math.Pi / float64(rays)
    drawFocusLine(screen, x, y, angle, alpha)
}`,
  },
  {
    slug:"vfx-faux-bloom",tier:"basic",step:"17",stars:"★★★☆☆",labKind:"bloom",
    concept:{ja:"拡大コピー + BlendLighter",en:"Scaled copies + BlendLighter"},hubDesc:{ja:"大きく薄い光を重ね、本物のぼかしなしで発光を作ります。",en:"Stack larger faint lights for a glow without a real blur shader."},
    ja:{navConcept:"疑似ブルーム",title:"大きな薄いコピーで疑似ブルーム",lead:"光の絵を少しずつ大きく、薄くして後ろへ重ねます。正確なぼかしではありませんが、中心から光があふれるブルームに見えます。シェーダーへ進む前に使える軽い方法です。",deepEyebrow:"OPTICAL TRICK / BLOOM",deepH:"ぼかさずに<br>ぼけた光を作る",deepLead:"中心は白く小さく、外側ほど大きく透明なコピーにします。加算合成で重なった中心が明るくなるため、目は輪郭のない光として読みます。暗い背景ほど効果的です。",concepts:[{h:"大きくする",p:"コピーごとに倍率を増やします。",code:"1 + i*spread"},{h:"薄くする",p:"外側は小さいアルファにします。",code:"ScaleAlpha(alpha)"},{h:"足し合わせる",p:"中心の明るさを積み上げます。",code:"BlendLighter"}],lab:{eyebrow:"TRY IT / BLOOM",title:"光の広がりを比べよう",p:"コピー枚数と広がりを変え、芯と光のにじみを見比べます。"},codeHead:{eyebrow:"IN THE WASM",h:"大きな半透明を後ろへ",p:"ぼかし処理をせず、同じ柔らかい円を数回重ねています。"},whys:[{eyebrow:"WHY IT WORKS",h:"輪郭が重なって消える",p:"透明な層の連続を、目が滑らかな光として補います。"},{eyebrow:"WHY ADDITIVE",h:"中心だけ強くなる",p:"重なりの多い中心が白へ近づきます。"},{eyebrow:"TRY NEXT",h:"本物のBlurへ",p:"オフスクリーン画像とKageで、次は画面全体をぼかせます。"}]},
    en:{navConcept:"Faux bloom",title:"Faux Bloom from Larger Copies",lead:"Draw the same light larger and fainter behind itself. It is not a true blur, but the eye reads a bright center spilling into bloom—a cheap step before shaders.",deepEyebrow:"OPTICAL TRICK / BLOOM",deepH:"Make blurred light<br>without blurring",deepLead:"Keep the center small and white, with larger transparent copies outside. Additive overlap brightens the middle, so the eye reads edgeless light. It works best on dark backgrounds.",concepts:[{h:"Enlarge",p:"Increase scale for each copy.",code:"1 + i*spread"},{h:"Fade",p:"Use low alpha outside.",code:"ScaleAlpha(alpha)"},{h:"Add",p:"Pile brightness in the center.",code:"BlendLighter"}],lab:{eyebrow:"TRY IT / BLOOM",title:"Compare the glow spread",p:"Change copy count and spread, then compare the core with its spill."},codeHead:{eyebrow:"IN THE WASM",h:"Large translucent copies behind",p:"No blur pass—just several draws of the same soft circle."},whys:[{eyebrow:"WHY IT WORKS",h:"Edges merge",p:"The eye fills transparent layers into smooth light."},{eyebrow:"WHY ADDITIVE",h:"The center gets strongest",p:"More overlap pushes the core toward white."},{eyebrow:"TRY NEXT",h:"Move to real blur",p:"Use an offscreen image and Kage to blur a whole layer."}]},
    code:`for i := 0; i < copies; i++ {
    scale := 1 + float64(i)*spread
    op.GeoM.Scale(scale, scale)
    op.ColorScale.ScaleAlpha(alpha)
    op.Blend = ebiten.BlendLighter
    screen.DrawImage(glow, op)
}`,
  },
].map((lesson) => ({
  ...lesson,
  ja: {...lesson.ja, deepEyebrow: `DEEP DIVE / ${lesson.ja.deepEyebrow}`},
  en: {...lesson.en, deepEyebrow: `DEEP DIVE / ${lesson.en.deepEyebrow}`},
}));
