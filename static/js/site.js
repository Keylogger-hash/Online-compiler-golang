// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/**
 * A bit of navigation related code for handling dismissible elements.
 */
window.initFuncs = [];

(() => {
  'use strict';

  function registerHeaderListeners() {
    const header = document.querySelector('.js-header');
    const menuButtons = document.querySelectorAll('.js-headerMenuButton');
    menuButtons.forEach(button => {
      button.addEventListener('click', e => {
        e.preventDefault();
        header.classList.toggle('is-active');
        button.setAttribute(
          'aria-expanded',
          header.classList.contains('is-active')
        );
      });
    });

    const scrim = document.querySelector('.js-scrim');
    scrim.addEventListener('click', e => {
      e.preventDefault();
      header.classList.remove('is-active');
      menuButtons.forEach(button => {
        button.setAttribute(
          'aria-expanded',
          header.classList.contains('is-active')
        );
      });
    });
  }

  function registerSolutionsTabs() {
    // Handle tab navigation on Solutions page.
    const tabList = document.querySelector('.js-solutionsTabs');

    if (tabList) {
      const tabs = tabList.querySelectorAll('[role="tab"]');
      let tabFocus = getTabFocus();

      changeTabs({ target: tabs[tabFocus] })

      tabs.forEach(tab => {
        tab.addEventListener('click', changeTabs);
      });

      // Enable arrow navigation between tabs in the tab list
      tabList.addEventListener('keydown', e => {
        // Move right
        if (e.keyCode === 39 || e.keyCode === 37) {
          tabs[tabFocus].setAttribute('tabindex', -1);
          if (e.keyCode === 39) {
            tabFocus++;
            // If we're at the end, go to the start
            if (tabFocus >= tabs.length) {
              tabFocus = 0;
            }
            // Move left
          } else if (e.keyCode === 37) {
            tabFocus--;
            // If we're at the start, move to the end
            if (tabFocus < 0) {
              tabFocus = tabs.length - 1;
            }
          }
          tabs[tabFocus].setAttribute('tabindex', 0);
          tabs[tabFocus].focus();
          setTabFocus(tabs[tabFocus].id);
        }
      });

      function getTabFocus() {
        const hash = window.location.hash;

        switch (hash) {
          case '#use-cases':
            return 1;
          case '#case-studies':
          default:
            return 0;
        }
      }

      function setTabFocus(id) {
        switch (id) {
          case 'btn-tech':
            tabFocus = 1;
            window.location.hash = '#use-cases';
            break;
          case 'btn-companies':
          default:
            window.location.hash = '#case-studies';
            tabFocus = 0;
        }
      }

      function changeTabs(e) {
        const target = e.target;
        const parent = target.parentNode;
        const grandparent = parent.parentNode;

        // Remove all current selected tabs
        parent
          .querySelectorAll('[aria-selected="true"]')
          .forEach(t => t.setAttribute('aria-selected', false));

        // Set this tab as selected
        target.setAttribute('aria-selected', true);
        setTabFocus(target.id)

        // Hide all tab panels
        grandparent
          .querySelectorAll('[role="tabpanel"]')
          .forEach(panel => panel.setAttribute('hidden', true));

        // Show the selected panel
        grandparent.parentNode
          .querySelector(`#${target.getAttribute('aria-controls')}`)
          .removeAttribute('hidden');
      }
    }
  }

  /**
   * Attempts to detect user's operating system and sets the download
   * links accordingly
   */
  async function setDownloadLinks() {
    const versionElement = document.querySelector('.js-latestGoVersion');
    if (versionElement) {
      const anchorTagWindows = document.querySelector('.js-downloadWin');
      const anchorTagMac = document.querySelector('.js-downloadMac');
      const anchorTagLinux = document.querySelector('.js-downloadLinux');
      const version = await getLatestVersion();

      const macDownloadUrl = `/dl/${version}.darwin-amd64.pkg`;
      const windowsDownloadUrl = `/dl/${version}.windows-amd64.msi`;
      const linuxDownloadUrl = `/dl/${version}.linux-amd64.tar.gz`;

      anchorTagWindows.href = windowsDownloadUrl;
      anchorTagMac.href = macDownloadUrl;
      anchorTagLinux.href = linuxDownloadUrl;

      /*
       * Note: we do not change .js-downloadBtn anymore
       * because it is impossible to tell reliably which architecture
       * the user's browser is running on.
       */
    }
  }

  /**
   * Retrieves list of Go versions & returns the latest
   */
  async function getLatestVersion() {
    let version = 'go1.17'; // fallback version if fetch fails
    try {
      const versionData = await (
        await fetch('/dl/?mode=json')
      ).json();
      if (!versionData.length) {
        return version;
      }
      versionData.sort((v1, v2) => {
        return v2.version - v1.version;
      });
      version = versionData[0].version;
    } catch (err) {
      console.error(err);
    }
    return version;
  }

  /**
   * initialThemeSetup sets data-theme attribute based on preferred color
   */

  function initialThemeSetup() {
    const theme = document.cookie.match(/prefers-color-scheme=(light|dark|auto)/)?.[1]

    if (theme) {
      document.querySelector('html').setAttribute('data-theme', theme);
    }
  }

  /**
   * setThemeButtons sets click listeners for toggling theme buttons
   */
  function setThemeButtons() {
    for (const el of document.querySelectorAll('.js-toggleTheme')) {
      el.addEventListener('click', () => {
        toggleTheme();
      });
    }
  }

  /**
   * toggleTheme switches the preferred color scheme between auto, light, and dark.
   */
  function toggleTheme() {
    let nextTheme = 'dark';
    const theme = document.documentElement.getAttribute('data-theme');
    if (theme === 'dark') {
      nextTheme = 'light';
    } else if (theme === 'light') {
      nextTheme = 'auto';
    }
    let domain = '';
    if (location.hostname === 'go.dev') {
      // Include subdomains to apply the setting to pkg.go.dev.
      domain = 'domain=.go.dev;';
    }
    document.documentElement.setAttribute('data-theme', nextTheme);
    document.cookie =
      `prefers-color-scheme=${nextTheme};${domain}path=/;max-age=31536000;`;
  }

  initialThemeSetup();

  window.addEventListener('DOMContentLoaded', () => {
    registerHeaderListeners();
    registerSolutionsTabs();
    setDownloadLinks();
    setThemeButtons();
  });
})();
