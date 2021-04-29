#!/bin/sh
CURDIR=`pwd`

make amd64
make install-amd64 DESTDIR=${CURDIR}/release/build/
cd ${CURDIR}
find release/build -type f | xargs md5sum | sed 's/  release\/build\// /' > ${CURDIR}/dist/DEBIAN.amd64/md5sums
chmod a+x ${CURDIR}/dist/DEBIAN.amd64/postinst
chmod a+x ${CURDIR}/dist/DEBIAN.amd64/postrm
chmod a+x ${CURDIR}/dist/DEBIAN.amd64/prerm
rm -rf ${CURDIR}/release/build/DEBIAN
cp -rf ${CURDIR}/dist/DEBIAN.amd64 ${CURDIR}/release/build/DEBIAN
cp -rf ${CURDIR}/dist/etc ${CURDIR}/release/build
dpkg-deb -b ${CURDIR}/release/build ${CURDIR}/release