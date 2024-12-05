ARG VARIANT=stable-slim
FROM docker.io/library/debian:$VARIANT

ENV DEBIAN_FRONTEND=noninteractive

RUN \
    apt-get update; \
    apt-get install --no-install-recommends --assume-yes \
    biber \
    chktex \
    cm-super \
    dvidvi \
    dvipng \
    fragmaster \
    git \
    lacheck \
    latexdiff \
    latexmk \
    lcdf-typetools \
    lmodern \
    make \
    psutils \
    purifyeps \
    t1utils \
    tex-gyre \
    texinfo \
    texlive-base \
    texlive-bibtex-extra \
    texlive-binaries \
    texlive-extra-utils \
    texlive-font-utils \
    texlive-fonts-extra \
    texlive-fonts-extra-links \
    texlive-fonts-recommended \
    texlive-formats-extra \
    texlive-lang-english \
    texlive-lang-european \
    texlive-latex-base \
    texlive-latex-extra \
    texlive-latex-recommended \
    texlive-luatex \
    texlive-metapost \
    texlive-pictures \
    texlive-plain-generic \
    texlive-pstricks \
    texlive-science \
    texlive-xetex; \
    apt-get clean && rm -fr /var/lib/apt/lists/*
