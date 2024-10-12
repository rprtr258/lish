#!/usr/bin/env python
import os
import argparse
from collections import Counter
from itertools import islice

def first_n_sorted_dict_items(d, to_take=10):
    return islice(sorted(d.items(), key=lambda kv:-kv[1]), to_take)

def format_line(cnt, overall, word):
    return f"{cnt:3d} ({100*cnt/overall:2.0f}%) {word}"

parser = argparse.ArgumentParser(description="Print most frequent called commands.")
parser.add_argument(
    "-n",
    "--numcmds",
    dest="numcmds",
    type=int,
    default=10,
    help="how many most frequent commands to print (default: 10)"
)
parser.add_argument(
    "-m",
    "--numlines",
    dest="numlines",
    type=int,
    default=3,
    help="how many most frequent lines to print (default: 3)"
)
parser.add_argument(
    "-s",
    "--source",
    dest="file",
    type=str,
    default=os.path.join(os.environ["HOME"], ".bash_history"),
    help="bash history file (default: ~/.bash_history)"
)

args = parser.parse_args()
history_file = args.file
lines = []

with open(history_file) as fd:
    for line in fd:
        if len(line.split()) == 0:
            continue
        lines.append(line.strip())

words_count = Counter(w.split()[0] for w in lines)
all_words = sum(words_count.values())
for word, count in first_n_sorted_dict_items(words_count, args.numcmds):
    print(format_line(count, all_words, word))
    cmds = [line for line in lines if word in line.split()]
    for cmd, cmd_count in first_n_sorted_dict_items(Counter(cmds), args.numlines):
        if len(cmd) > 50:
            cmd = cmd[:25] + "..." + cmd[-25:]
        print(f"  " + format_line(cmd_count, count, cmd))

