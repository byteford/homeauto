FROM homeassistant/home-assistant

#COPY configuration.yaml /config/configuration.yaml
COPY ./config/ /config
WORKDIR /config
RUN ls && wget -q -O - https://hacs.xyz/install | bash -
