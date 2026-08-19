# -*- coding: utf-8 -*-
import json

PASSAGES = [
  {
    "title": "coffee",
    "passage": ("The history of coffee begins in the highlands of Ethiopia, where legend says a goat herder "
      "named Kaldi noticed his goats became unusually energetic after eating berries from a certain tree. "
      "The berries were roasted, ground, and brewed, and the drink spread across the Arabian Peninsula. "
      "By the 1500s, coffeehouses had appeared in cities such as Mecca and Cairo, serving as places for "
      "conversation, music, and business. Coffee reached Europe in the seventeenth century and quickly "
      "became a popular beverage. Some rulers distrusted coffeehouses, fearing they encouraged political "
      "discussion. Nevertheless, the drink's popularity continued to grow, and coffee cultivation expanded "
      "to colonies in Asia and the Americas."),
    "questions": [
      ("What is the main idea of the passage?",
       ["Coffee spread from Ethiopia to become a popular drink around the world.",
        "Kaldi was the first person to grow coffee in Africa.",
        "Coffeehouses were forbidden in every country in Europe.",
        "Coffee cultivation began in the Americas in the 1500s."], 0, "easy",
       "A. Ide utamanya perjalanan kopi: dari Ethiopia menyebar ke Jazirah Arab, Eropa, sampai Amerika. B, C, dan D terlalu spesifik atau bertentangan dengan isi teks, misalnya Kaldi menemukan bukan menanam, dan budidaya di Amerika terjadi belakangan."),
      ("According to the passage, what did Kaldi notice about his goats?",
       ["They became unusually energetic after eating certain berries.",
        "They refused to eat the berries from the tree.",
        "They slept longer after drinking the brew.",
        "They led him to a coffeehouse in Mecca."], 0, "easy",
       "A. Paragraf pertama jelas bilang kambingnya jadi 'unusually energetic' setelah makan buah beri. B, C, dan D nggak ada di teks."),
      ("Why did some rulers distrust coffeehouses?",
       ["They feared the places encouraged political discussion.",
        "The coffee served there was too expensive.",
        "Coffeehouses attracted too many farmers.",
        "The buildings were in poor condition."], 0, "medium",
       "A. Teks bilang 'Some rulers distrusted coffeehouses, fearing they encouraged political discussion'. B, C, dan D nggak pernah disebut."),
      ("The word 'beverage' in the passage is closest in meaning to:",
       ["a type of food", "a drink", "a container", "a ceremony"], 1, "easy",
       "B. a drink. 'Beverage' = minuman. Kopi disebut 'popular beverage' karena menjadi minuman favorit. 'Food', 'container', dan 'ceremony' nggak nyambung."),
      ("It can be inferred from the passage that coffee cultivation spread to colonies because ...",
       ["coffee was popular and therefore profitable to grow.",
        "coffeehouses were banned in Europe.",
        "goats ate all the coffee berries.",
        "the Arabian Peninsula had no farmland."], 0, "hard",
       "A. Kalau kopi populer, otomatis permintaannya besar, jadi wajar dibudidayakan di koloni. Itu kesimpulan logis dari 'popularity continued to grow'. B, C, dan D nggak didukung teks."),
    ],
  },
  {
    "title": "reef",
    "passage": ("The Great Barrier Reef, stretching more than 2,300 kilometers along the coast of northeastern "
      "Australia, is the largest coral reef system on Earth. It is home to thousands of species of fish, "
      "mollusks, turtles, and seabirds. The reef is built by tiny animals called coral polyps, which produce "
      "a hard skeleton of calcium carbonate. These skeletons accumulate over centuries to form the reef "
      "structure. Coral polyps live in a partnership with microscopic algae that provide them with energy "
      "through photosynthesis. When ocean temperatures rise, the algae may leave the coral, causing it to "
      "lose its color and turn white, a process known as bleaching. Although coral can recover, severe or "
      "repeated bleaching can kill entire sections of the reef. Scientists warn that climate change is the "
      "greatest threat to the reef's future."),
    "questions": [
      ("What is the main idea of the passage?",
       ["The Great Barrier Reef is threatened mainly by rising ocean temperatures.",
        "Coral reefs are built only by seabirds and turtles.",
        "The Great Barrier Reef is the smallest reef system in the world.",
        "Scientists have discovered how to stop all coral bleaching."], 0, "easy",
       "A. Teks menjelaskan apa itu terumbu karang, lalu inti di akhir: kenaikan suhu dan perubahan iklim jadi ancaman terbesar. B, C, dan D bertentangan dengan teks."),
      ("What causes coral to turn white?",
       ["The algae leave the coral when temperatures rise.",
        "The coral produces extra calcium carbonate.",
        "Fish remove the algae from the reef.",
        "The reef moves to warmer water."], 0, "medium",
       "A. Prosesnya: suhu naik, alga pergi, karang kehilangan warna jadi putih (bleaching). B, C, dan D nggak sesuai dengan penjelasan di teks."),
      ("The word 'accumulate' in the passage is closest in meaning to:",
       ["to build up over time", "to disappear quickly", "to break apart suddenly", "to move from place to place"], 0, "medium",
       "A. to build up over time. 'Accumulate' = menumpuk atau terbentuk sedikit demi sedikit selama berabad-abad ('over centuries'). B, C, dan D kebalikan atau beda makna."),
      ("According to the passage, what do coral polyps produce?",
       ["a hard skeleton of calcium carbonate", "a supply of microscopic algae", "large amounts of fresh water", "a thick layer of sand"], 0, "easy",
       "A. Teks bilang polip menghasilkan 'a hard skeleton of calcium carbonate'. Alga justru penyedia energi untuk polip, bukan produksi polip. B, C, dan D salah."),
      ("The passage suggests that the reef's future depends mainly on ...",
       ["controlling the rise in ocean temperatures.",
        "reducing the number of seabirds.",
        "removing the calcium carbonate skeletons.",
        "increasing tourism along the coast."], 0, "hard",
       "A. Kalau ancaman terbesar adalah perubahan iklim, masa depan terumbu bergantung pada mengendalikan kenaikan suhu. Itu kesimpulan logis dari kalimat penutup. Opsi lain nggak nyambung."),
    ],
  },
]

Q = []
for p in PASSAGES:
    for i, (qt, opts, ci, diff, exp) in enumerate(p["questions"]):
        Q.append({
            "section": "reading", "type": "reading-comprehension", "difficulty": diff,
            "passage": p["passage"],
            "question_text": qt,
            "options": opts, "correct_index": ci,
            "explanation": exp,
            "highlights": {},
        })

with open("/home/nrizq/Documents/Codes/toefl-prep/generated/reading_a.json", "w") as f:
    json.dump(Q, f, ensure_ascii=False, indent=1)
print("reading_a:", len(Q))