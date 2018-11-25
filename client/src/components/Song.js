import React from 'react';
import styles from './Song.module.css';

const Song = ({name, subName, authorName, image}) => (
  <div className={styles.song}>
    <div className={styles.image} />
    <div className={styles.name}>{name}</div>
    <div className={styles.subName}>{subName}</div>
    <div className={styles.author}>{authorName}</div>
  </div>
);

export default Song;
