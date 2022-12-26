cd src
rm ahfs2vec
clear
echo building
go build
echo 'running ahfs\n'
./ahfs2vec
cd ..
dot -Tpng vis/g.dot > vis/graph.png
dot -Tsvg vis/g.dot > vis/graph.svg

echo "****"
echo "some meds are duplicated in the graph. (ex, cetrizine)"
echo "****"

python ahfs2vec.py --dims 64 --force-train true