# Gitea 메일 템플릿

[Gitea](https://about.gitea.com) 인스턴스를 위한 메일 템플릿 컬렉션.

> **110 파일 — 10 스타일, 각 11종류**

---

## 스타일 갤러리

| 미리보기 | 스타일 | 대상 | 특징 |
|---|---|---|---|
| ![Horizon](images/horizon.png) | **Horizon** | 엔터프라이즈 | Blue accent, slate, centred cards |
| ![Terminal](images/terminal.png) | **Terminal** | 개발자 | Dark mode, monospace, green CLI |
| ![Ember](images/ember.png) | **Ember** | 커뮤니티 | Warm amber, rounded, inclusive |
| ![Bloom](images/bloom.png) | **Bloom** | 크리에이티브 | Frosted glass, soft blue light |
| ![Heritage](images/heritage.png) | **Heritage** | 교육/연구 | Navy and gold, serif, classic |
| ![Neon](images/neon.png) | **Neon** | 게임/Web3 | Cyberpunk neon, pink and cyan |
| ![Mono](images/mono.png) | **Mono** | 디자인/편집 | Swiss brutalist, black-white-red |
| ![Terra](images/terra.png) | **Terra** | 지속가능성 | Warm earth tones, organic textures |
| ![Ink](images/ink.png) | **Ink** | 출판/뉴스 | Editorial print, navy and gold |
| ![Aurora](images/aurora.png) | **Aurora** | 프리미엄SaaS | Ethereal gradients, purple and teal |

> 이미지는 [라이브 프리뷰](../preview/index.html)에서 캡처한 600px 너비 스크린샷입니다. 캡처 방법은 [images/README.md](images/README.md)를 참조하세요.

[**라이브 프리뷰**](../preview/index.html)

---

## 설치

```bash
cp -r themes/horizon/mail/* /var/lib/gitea/custom/templates/mail/
systemctl restart gitea
```

## 프리뷰

**정적 모드:**
```bash
cd tools && go run . preview all
```
그런 다음 `preview/index.html` 열기.

**개발 서버 (실시간 리로드 + Juice CSS 인라인 + 클라이언트 시뮬레이션):**
```bash
cd tools && go run . dev     # Node.js 필요
# → http://localhost:3456
```

> [!NOTE]
> Dev 시뮬레이션은 각 이메일 클라이언트의 렌더링을 완전히 재현할 수 없습니다 — 참고용입니다. 실제 클라이언트에서 반드시 확인하세요.

## 호환성

- **Gitea 1.21+**, 100% 호환, 내장 함수만 사용

## 라이선스

MIT — [LICENSE](../LICENSE).
