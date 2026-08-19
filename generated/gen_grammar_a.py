# -*- coding: utf-8 -*-
import json

Q = [
  ("Neither the director nor the committee members was informed about the change.",
   ["Neither the director", "nor the committee members was", "informed about", "the change"], 1, "medium",
   "B. Aturan 'neither...nor': verb mengikuti subjek terdekat, yaitu 'the committee members' (jamak), jadi 'were informed', bukan 'was'. Agreement jamak-tunggal ini jebakan khas TOEFL.",
   {"conjunction": ["Neither", "nor"], "noun": ["director", "committee", "members", "change"], "verb": ["informed"]}),
  ("The quality of the products have improved considerably this year.",
   ["The quality", "of the products have", "improved considerably", "this year"], 1, "medium",
   "B. Subjek utamanya 'the quality' (tunggal); 'of the products' cuma modifier, jadi 'has improved', bukan 'have'. Jebakannya: 'products' jamak bikin kelihatan harus 'have'.",
   {"noun": ["quality", "products", "year"], "verb": ["improved"]}),
  ("She enjoys reading, to swim, and hiking on the weekends.",
   ["She enjoys", "reading, to swim,", "and hiking", "on the weekends"], 1, "medium",
   "B. Paralelisme: 'reading' dan 'hiking' gerund, jadi yang tengah harus 'swimming', bukan 'to swim'. Rangkaian harus sejajar bentuknya.",
   {"verb": ["enjoys", "reading", "to swim", "hiking"]}),
  ("The training teaches employees to communicate clearly, working in teams, and manage conflicts.",
   ["The training teaches employees", "to communicate clearly,", "working in teams,", "and manage conflicts"], 2, "hard",
   "C. Paralelisme: 'to communicate' dan 'to manage' infinitive, jadi yang tengah harus 'to work in teams', bukan 'working'. Ketiganya harus sejajar.",
   {"noun": ["training", "employees", "teams", "conflicts"], "verb": ["communicate", "working", "manage"]}),
  ("Her and I are responsible for organizing the event.",
   ["Her and I", "are responsible", "for organizing", "the event"], 0, "easy",
   "A. 'Her and I' berfungsi sebagai subjek, jadi harus bentuk subjek 'She and I'. 'Her' adalah objek pronoun dan nggak bisa jadi subjek.",
   {"pronoun": ["Her", "I"], "verb": ["responsible", "organizing"], "noun": ["event"]}),
  ("Every participant must bring their own laptop to the training.",
   ["Every participant must bring", "their own laptop", "to the training", "(no error)"], 1, "hard",
   "B. Aturan klasik TOEFL: 'every participant' tunggal, jadi pronoun harus 'his or her', bukan 'their'. Ini jebakan pronoun agreement yang masih sering diujikan.",
   {"determiner": ["Every"], "noun": ["participant", "laptop", "training"], "pronoun": ["their"]}),
  ("By the time we arrived at the theater, the play has already begun.",
   ["By the time we arrived", "at the theater", "the play has already begun", "(no error)"], 2, "hard",
   "C. Sequence of tenses: 'arrived' (past) dan 'begun' terjadi lebih dulu, jadi harus past perfect 'had already begun'. 'Has' (present perfect) nggak cocok dengan konteks masa lalu.",
   {"noun": ["time", "theater", "play"], "verb": ["arrived", "begun"]}),
  ("She said that she will call me as soon as she arrives.",
   ["She said that", "she will call me", "as soon as she arrives", "(no error)"], 1, "medium",
   "B. Reported speech: kalau pembicaranya lampau ('said'), verb-nya mundur jadi 'would call', bukan 'will call'. Ini aturan tense agreement dalam kalimat tidak langsung.",
   {"verb": ["said", "will call", "arrives"]}),
  ("The new model is more better than the previous one.",
   ["The new model", "is more better", "than the previous one", "(no error)"], 1, "easy",
   "B. Double comparative: 'better' sudah bentuk comparative, nggak boleh ditambah 'more'. Jadi 'is better', bukan 'is more better'.",
   {"adjective": ["better"], "noun": ["model"], "pronoun": ["one"]}),
  ("He is the most friendliest person I have ever met.",
   ["He is the most", "friendliest person", "I have ever met", "(no error)"], 0, "medium",
   "A. Superlative ganda: 'friendliest' sudah bentuk superlative, nggak perlu 'the most'. Cukup 'the friendliest person'.",
   {"adjective": ["most", "friendliest"], "noun": ["person"], "verb": ["met"]}),
]

out = []
for text, segs, ci, diff, exp, hl in Q:
    out.append({
        "section": "grammar_adv", "type": "error-identification", "difficulty": diff,
        "question_text": text,
        "options": segs, "correct_index": ci,
        "explanation": exp,
        "highlights": hl,
    })

with open("/home/nrizq/Documents/Codes/toefl-prep/generated/grammar_a.json", "w") as f:
    json.dump(out, f, ensure_ascii=False, indent=1)
print("grammar_a:", len(out))