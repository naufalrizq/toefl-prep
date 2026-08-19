# -*- coding: utf-8 -*-
import json

VOCAB = [
  ("mitigate", "verb", ["to worsen", "to lessen", "to ignore", "to emphasize"], 1, "hard",
   "B. to lessen. 'Mitigate' artinya mengurangi dampak buruk. 'To worsen' kebalikannya, sedangkan 'ignore' dan 'emphasize' nggak ada hubungannya dengan mengurangi."),
  ("obsolete", "adjective", ["modern", "out of date", "expensive", "useful"], 1, "easy",
   "B. out of date. 'Obsolete' = nggak dipakai lagi karena sudah ketinggalan zaman. 'Modern' justru kebalikannya."),
  ("ubiquitous", "adjective", ["rare", "found everywhere", "unexpected", "harmful"], 1, "medium",
   "B. found everywhere. 'Ubiquitous' = ada di mana-mana. 'Rare' kebalikannya, 'unexpected' dan 'harmful' nggak nyambung."),
  ("ambiguous", "adjective", ["clear", "unclear", "brief", "urgent"], 1, "medium",
   "B. unclear. 'Ambiguous' = maknanya ganda, bisa ditafsirkan beda-beda, jadi nggak jelas. 'Clear' kebalikannya."),
  ("concise", "adjective", ["brief", "wordy", "vague", "lengthy"], 0, "easy",
   "A. brief. 'Concise' = ringkas dan padat, nggak bertele-tele. 'Wordy' dan 'lengthy' kebalikannya."),
  ("detrimental", "adjective", ["helpful", "harmful", "minor", "temporary"], 1, "medium",
   "B. harmful. 'Detrimental' = merugikan atau merusak. 'Helpful' kebalikannya."),
  ("diverse", "adjective", ["varied", "identical", "limited", "simple"], 0, "easy",
   "A. varied. 'Diverse' = beragam, macam-macam. 'Identical' dan 'limited' kebalikannya."),
  ("eliminate", "verb", ["to remove", "to add", "to reduce", "to avoid"], 0, "easy",
   "A. to remove. 'Eliminate' = menyingkirkan atau membuang sepenuhnya. 'To add' kebalikannya."),
  ("enhance", "verb", ["to improve", "to worsen", "to replace", "to delay"], 0, "medium",
   "A. to improve. 'Enhance' = meningkatkan atau memperbaiki. 'To worsen' kebalikannya."),
  ("feasible", "adjective", ["impossible", "practical", "expensive", "urgent"], 1, "medium",
   "B. practical. 'Feasible' = bisa dilaksanakan, masuk akal secara praktis. 'Impossible' kebalikannya."),
  ("genuine", "adjective", ["authentic", "imitation", "faulty", "common"], 0, "easy",
   "A. authentic. 'Genuine' = asli, bukan tiruan. 'Imitation' justru artinya tiruan."),
  ("hypothetical", "adjective", ["assumed", "proven", "actual", "repeated"], 0, "hard",
   "A. assumed. 'Hypothetical' = bersifat dugaan atau teori, belum terbukti. 'Proven' dan 'actual' kebalikannya."),
  ("inevitable", "adjective", ["unavoidable", "preventable", "optional", "sudden"], 0, "hard",
   "A. unavoidable. 'Inevitable' = nggak bisa dihindari, pasti terjadi. 'Preventable' dan 'optional' kebalikannya."),
]

Q = []
for word, pos, opts, ci, diff, exp in VOCAB:
    q = {
        "section": "vocabulary", "type": "vocab-multiple-choice", "difficulty": diff,
        "question_text": "The word '" + word + "' is closest in meaning to:",
        "options": opts, "correct_index": ci,
        "explanation": exp,
        "highlights": {pos: [word]},
    }
    Q.append(q)

with open("/home/nrizq/Documents/Codes/toefl-prep/generated/vocab_a.json", "w") as f:
    json.dump(Q, f, ensure_ascii=False, indent=1)
print("vocab_a:", len(Q))