et -e # exit with nonzero exit code if anything fails

# go to the out directory and create a *new* Git repo
cd build
git init

# inside this git repo we'll pretend to be a new user
git config user.name "Auto Coverage"
git config user.email "jsuz.iicot.o@gmail.com"

# The first and only commit to this new Git repo contains all the
# files present with the commit message "Deploy to GitHub Pages".
git add .
git commit -m "Deploy to GitHub Pages"

# Force push from the current repo's master branch to the remote
# repo's gh-pages branch. (All previous history on the gh-pages branch
# will be lost, since we are overwriting it.) We redirect any output to
# /dev/null to hide any sensitive credential data that might otherwise be exposed.
echo 'Pushing github pages'
git push --force --quiet "https://${GH_TOKEN}@github.com/suzujun/gendao.git" master:gh-pages
