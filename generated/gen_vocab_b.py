# -*- coding: utf-8 -*-
import json

VOCAB = [
  ("notorious", "adjective", ["widely known for something bad", "completely unknown", "highly respected", "slightly famous"], 0, "hard",
   "A. widely known for something bad. 'Notorious' = terkenal karena hal buruk. 'Highly respected' kebalikannya."),
  ("pragmatic", "adjective", ["practical", "idealistic", "careless", "lengthy"], 0, "medium",
   "A. practical. 'Pragmatic' = realistis, ngandelin hal praktis daripada teori. 'Idealistic' kebalikannya."),
  ("reluctant", "adjective", ["unwilling", "eager", "angry", "confused"], 0, "medium",
   "A. unwilling. 'Reluctant' = ogah-ogahan, enggan. 'Eager' = antusias, kebalikannya."),
  ("resilient", "adjective", ["able to recover quickly", "easily damaged", "permanently broken", "slow to react"], 0, "hard",
   "A. able to recover quickly. 'Resilient' = tahan banting, cepat pulih. 'Easily damaged' kebalikannya."),
  ("scrutinize", "verb", ["to examine closely", "to ignore completely", "to improve quickly", "to repeat loudly"], 0, "hard",
   "A. to examine closely. 'Scrutinize' = memeriksa dengan teliti banget. 'To ignore completely' kebalikannya."),
  ("subsequent", "adjective", ["following", "previous", "simultaneous", "unrelated"], 0, "medium",
   "A. following. 'Subsequent' = yang menyusul, setelahnya. 'Previous' kebalikannya."),
  ("sufficient", "adjective", ["enough", "scarce", "excessive", "minimal"], 0, "easy",
   "A. enough. 'Sufficient' = cukup. 'Scarce' = langka atau kurang."),
  ("tedious", "adjective", ["boring", "exciting", "quick", "simple"], 0, "medium",
   "A. boring. 'Tedious' = membosankan dan melelahkan karena lama atau monoton. 'Exciting' kebalikannya."),
  ("undermine", "verb", ["to weaken", "to strengthen", "to support", "to explain"], 0, "hard",
   "A. to weaken. 'Undermine' = melemahkan dari dalam atau merusak perlahan. 'To strengthen' kebalikannya."),
  ("vital", "adjective", ["essential", "optional", "harmful", "minor"], 0, "easy",
   "A. essential. 'Vital' = sangat penting, nggak bisa diabaikan. 'Optional' dan 'minor' kebalikannya."),
  ("candid", "adjective", ["frank", "secretive", "polite", "vague"], 0, "medium",
   "A. frank. 'Candid' = terus terang, jujur apa adanya. 'Secretive' kebalikannya."),
  ("imperative", "adjective", ["absolutely necessary", "slightly helpful", "commonly avoided", "mildly optional"], 0, "hard",
   "A. absolutely necessary. 'Imperative' = sangat wajib, nggak bisa ditawar. 'Optional' kebalikannya."),
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

with open("/home/nrizq/Documents/Codes/toefl-prep/generated/vocab_b.json", "w") as f:
    json.dump(Q, f, ensure_ascii=False, indent=1)
print("vocab_b:", len(Q))