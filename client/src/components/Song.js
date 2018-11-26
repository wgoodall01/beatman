import React from 'react';
import styles from './Song.module.css';
import classnames from 'classnames';
import gql from 'graphql-tag';

const play = uri => {
  const audio = new Audio(uri);
  audio.play();
};

const Song = ({name, subName, authorName, coverUri, audioUri}) => (
  <div className={styles.song}>
    <div className={styles.image}>
      <img src={coverUri} />
    </div>
    <div className={classnames(styles.name, styles.wrapText)}>{name}</div>
    <div className={styles.subName}>{subName}</div>
    <div className={styles.authorName}>{authorName}</div>
    <div className={styles.buttonBox}>
      <button onClick={() => play(audioUri)}>preview</button>
    </div>
  </div>
);

Song.fragment = gql`
  fragment Song on Song {
    name
    subName
    authorName
    coverUri
    audioUri
  }
`;

export default Song;
