# -*- coding: utf-8 -*-
import json

PASSAGES = [
  {
    "title": "aviation",
    "passage": ("The first controlled flight in a powered aircraft was achieved by Orville and Wilbur Wright "
      "in 1903 on the windswept dunes of Kitty Hawk, North Carolina. The flight lasted only twelve seconds "
      "and covered about thirty-seven meters, but it proved that controlled, powered flight was possible. "
      "In the following years, the Wright brothers continued to refine their machines, improving stability "
      "and control. Their work inspired a wave of aviation pioneers across the world, each racing to build "
      "faster and safer aircraft. Within two decades, aircraft had evolved from fragile wood-and-fabric "
      "machines into metal airliners carrying passengers between continents. The rapid pace of change was "
      "driven by the demand for speed in both civilian and military use, and aviation soon transformed how "
      "people traveled, communicated, and fought."),
    "questions": [
      ("What is the main idea of the passage?",
       ["The Wright brothers' first flight began a rapid evolution of aviation.",
        "The Wright brothers invented the airplane after decades of military use.",
        "Kitty Hawk was chosen because it had a large airport.",
        "Powered flight remained impossible for twenty years."], 0, "easy",
       "A. Teks mulai dari penerbangan pertama Wright bersaudara lalu menceritakan evolusi pesawat yang cepat. B salah urutan, C nggak disebut, D kebalikan dari isi teks."),
      ("How long did the Wright brothers' first flight last?",
       ["twelve seconds", "thirty-seven minutes", "two hours", "twenty years"], 0, "easy",
       "A. Teks bilang jelas 'The flight lasted only twelve seconds'. Opsi lain angka yang salah."),
      ("The word 'refine' in the passage is closest in meaning to:",
       ["to improve gradually", "to abandon completely", "to describe carefully", "to sell quickly"], 0, "medium",
       "A. to improve gradually. 'Refine' di sini = menyempurnakan mesin mereka pelan-pelan. 'Abandon' kebalikannya, 'describe' dan 'sell' nggak nyambung."),
      ("According to the passage, what drove the rapid development of aircraft?",
       ["demand for speed in civilian and military use",
        "the shortage of wood and fabric",
        "a law requiring metal airliners",
        "competition among Kitty Hawk farmers"], 0, "medium",
       "A. Teks bilang 'the rapid pace of change was driven by the demand for speed in both civilian and military use'. B, C, dan D nggak disebut."),
      ("It can be inferred from the passage that the first flight was significant because it ...",
       ["proved that controlled, powered flight was possible.",
        "allowed passengers to fly between continents.",
        "was the fastest aircraft of its time.",
        "used metal in the aircraft's construction."], 0, "hard",
       "A. Kalimat kunci: 'it proved that controlled, powered flight was possible', itulah makna pentingnya. B baru terjadi belakangan, C dan D nggak sesuai dengan pesawat pertama yang dari kayu dan kain."),
    ],
  },
  {
    "title": "sleep",
    "passage": ("Sleep is far more than a period of rest; it plays an essential role in how memories are "
      "formed and preserved. During deep sleep, the brain replays the experiences of the day, strengthening "
      "the connections between neurons. This process, called consolidation, helps transform fragile "
      "short-term memories into durable long-term ones. Studies show that people who sleep after learning a "
      "new skill tend to perform better than those who stay awake. The amount of sleep matters as well: both "
      "too little and too much sleep can interfere with memory. Researchers also note that a short nap of "
      "about twenty minutes can refresh attention, while longer naps may produce grogginess. For students, "
      "the practical advice is simple: after studying, a good night's sleep may be more valuable than an "
      "extra hour of review."),
    "questions": [
      ("What is the main idea of the passage?",
       ["Sleep helps turn short-term memories into lasting ones.",
        "All people need at least nine hours of sleep.",
        "Napping is harmful to academic performance.",
        "Memories are stored only during the day."], 0, "easy",
       "A. Inti teks: tidur penting untuk konsolidasi memori, mengubah memori jangka pendek jadi jangka panjang. B, C, dan D bertentangan atau nggak ada di teks."),
      ("What happens during deep sleep?",
       ["The brain replays the day's experiences and strengthens connections.",
        "The brain deletes all short-term memories.",
        "Neurons stop working completely.",
        "The body produces large amounts of calcium."], 0, "medium",
       "A. Teks bilang 'the brain replays the experiences of the day, strengthening the connections between neurons'. B, C, dan D nggak sesuai."),
      ("The word 'consolidation' in the passage is closest in meaning to:",
       ["the process of making something more stable",
        "a period of extreme tiredness",
        "a method of avoiding sleep",
        "the loss of all memory"], 0, "medium",
       "A. the process of making something more stable. Konsolidasi = proses memperkuat dan menstabilkan memori. B, C, dan D kebalikan atau beda makna."),
      ("According to the passage, what can a short nap of twenty minutes do?",
       ["refresh attention", "guarantee better exam scores", "replace all night sleep", "damage long-term memory"], 0, "easy",
       "A. Teks bilang 'a short nap of about twenty minutes can refresh attention'. Opsi lain dilebih-lebihkan atau nggak disebut."),
      ("The passage most likely suggests that students should ...",
       ["sleep well after studying instead of staying up late to review.",
        "avoid napping under any circumstances.",
        "study only in the morning when the brain is fresh.",
        "keep their sleep to less than six hours."], 0, "hard",
       "A. Kalimat penutup: 'after studying, a good night's sleep may be more valuable than an extra hour of review', jadi sarannya tidur yang cukup. B, C, dan D bertentangan."),
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

with open("/home/nrizq/Documents/Codes/toefl-prep/generated/reading_b.json", "w") as f:
    json.dump(Q, f, ensure_ascii=False, indent=1)
print("reading_b:", len(Q))