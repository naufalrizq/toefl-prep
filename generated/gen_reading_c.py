# -*- coding: utf-8 -*-
import json

PASSAGES = [
  {
    "title": "industrial",
    "passage": ("The Industrial Revolution, which began in Britain around the mid-eighteenth century, "
      "transformed economies that had been based on agriculture and handcrafts. New machines powered by "
      "steam allowed goods to be produced faster and more cheaply than ever before. Factories grew up near "
      "coal mines and rivers, and large numbers of people moved from rural villages to industrial cities in "
      "search of work. This urbanization changed society in fundamental ways: family life, housing, and "
      "daily routines were reshaped around the rhythm of the factory. Working conditions in the early "
      "factories were often harsh, with long hours and low pay, which eventually led to labor laws and the "
      "rise of trade unions. The revolution also accelerated the spread of railroads and steamships, "
      "shrinking distances and opening new markets around the world."),
    "questions": [
      ("What is the main idea of the passage?",
       ["The Industrial Revolution reshaped economies, cities, and daily life.",
        "Steam machines were invented only after 1900.",
        "Factory workers in Britain were never organized.",
        "The Industrial Revolution began in North America."], 0, "easy",
       "A. Teks merangkum dampak revolusi industri: ekonomi, urbanisasi, kondisi kerja, dan transportasi. B, C, dan D bertentangan dengan isi."),
      ("What powered the new machines of the Industrial Revolution?",
       ["steam", "electricity", "solar energy", "animal labor"], 0, "easy",
       "A. Teks bilang 'New machines powered by steam'. Listrik, tenaga surya, dan tenaga hewan nggak disebut."),
      ("According to the passage, why did many people move to industrial cities?",
       ["to search for work in the new factories",
        "to escape the smoke of coal mines",
        "to start farms near the rivers",
        "to attend the new universities"], 0, "medium",
       "A. Teks bilang orang pindah 'in search of work'. B, C, dan D nggak disebut."),
      ("The word 'rural' in the passage is closest in meaning to:",
       ["relating to the countryside", "relating to factories", "relating to money", "relating to steam power"], 0, "medium",
       "A. relating to the countryside. 'Rural villages' = desa di pedesaan, lawan dari 'industrial cities' yang disebut sesudahnya. Opsi lain nggak nyambung."),
      ("It can be inferred from the passage that early factory workers ...",
       ["faced difficult conditions that later improved through laws and unions.",
        "received generous salaries from the first day.",
        "worked only in rural areas.",
        "invented the steam engine themselves."], 0, "hard",
       "A. Teks bilang kondisi awal keras ('harsh, with long hours and low pay') lalu 'eventually led to labor laws and the rise of trade unions', jadi kondisi membaik lewat hukum dan serikat. B kebalikan, C dan D salah."),
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

with open("/home/nrizq/Documents/Codes/toefl-prep/generated/reading_c.json", "w") as f:
    json.dump(Q, f, ensure_ascii=False, indent=1)
print("reading_c:", len(Q))