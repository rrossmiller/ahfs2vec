from cProfile import label
import networkx as nx
import csv, json, os
import pandas as pd
from node2vec import Node2Vec
from tqdm import trange


def write_tsv(dimensions, file_path="."):
    print("generating TSVs")
    # gen_meta_labels()

    with open("vis/ahfs2vec.kv") as fin:
        file = fin.readlines()

    file.pop(0)
    vecs = []
    labels = ["ID\tName\n"]

    labels_df = pd.read_csv("vis/labels.tsv, header=None", sep="\t")
    labels_df.columns = ["name", "id"]
    labels_df["id"] = labels_df["id"].astype(str)
    print(labels_df)
    for i in trange(len(file)):
        splits = file[i].split()
        label = splits[0]

        # if it's a parent class, the root or a mistake, skip
        # if label[-2:] == "00" or label == "1" or label == "\\n":
            # continue
        # if it's a parent class, the root or a mistake, skip
        if label[0] != "3" or label == "1" or label == "\\n":
            continue

        label = labels_df.loc[labels_df["id"] == label].iloc[0].to_numpy()
        labels.append(f"{label[1].strip()}\t{label[0].strip()}\n")
        vecs.append(splits[1:])

    with open(f"{file_path}/ahfs_vecs-{dimensions}.tsv", "w") as fout:
        csv.writer(fout, delimiter="\t").writerows(vecs)

    labels[-1] = labels[-1].replace("\n", "")
    print(labels[-1])
    with open(f"{file_path}/ahfs_labels.tsv", "w") as fout:
        fout.writelines(labels)


def ahfs2vec(dimensions, force_train=False, workers=6):
    if not os.path.exists("vis/ahfs2vec.kv") or force_train:
        print("running ahfs2vec")
        g = nx.nx_pydot.read_dot("vis/g.dot")
        print("graph loaded")

        node2vec = Node2Vec(
            g, dimensions=dimensions, walk_length=30, num_walks=20, workers=workers
        )
        # Any keywords acceptable by gensim.Word2Vec can be passed, `dimensions` and `workers` are automatically passed (from the Node2Vec constructor)
        model = node2vec.fit(window=10, min_count=1, batch_words=4)

        # Save embeddings for later use
        model.wv.save_word2vec_format("vis/ahfs2vec.kv")
    else:
        print("vis/ahfs2vec.kv already exists")


if __name__ == "__main__":
    import argparse

    # os.system("clear")
    parser = argparse.ArgumentParser()
    parser.add_argument("--dims", default=8)
    parser.add_argument("--workers", default=6)
    parser.add_argument("--force-train", default=False)
    parser.add_argument("--path", default=".")

    dims = parser.parse_args().dims
    dims = int(dims)
    workers = int(parser.parse_args().workers)
    force_train = bool(parser.parse_args().force_train)
    path = parser.parse_args().path

    print(f"dims= {dims }")
    print(f"workers= {workers }")
    print(f"force_train= {force_train}")

    ahfs2vec(dims, workers=workers, force_train=force_train)
    print("ahfs2vec done")
    write_tsv(dims, file_path=path)
