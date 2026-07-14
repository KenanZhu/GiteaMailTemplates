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

### 작동 확인

관리자 테스트 이메일은 커스텀 템플릿을 사용하지 않습니다. 템플릿이 활성화
되었는지 확인하려면 실제 이메일 알림을 트리거하세요. 가장 빠른 방법은 비밀번호
재설정입니다: 로그아웃 후 로그인 페이지에서 **"Forgot password"** 를 클릭하고
재설정 이메일을 확인하세요.

## 프리뷰

**정적 모드:**
```bash
cd tools && go run . preview all
```
그런 다음 `preview/index.html` 열기.

**개발 서버 (실시간 리로드):**
```bash
cd tools && go run . dev 
# → http://localhost:3456
```

## 호환성

- **Gitea 1.25+** — v1.25에서 도입된 메일 템플릿 디렉토리 구조 사용
- **최신 테스트:** Gitea 1.27.0<!-- TRACKER:LATEST-TESTED -->
- Gitea 공식 템플릿과 100% 호환 — 자세한 내용은 [COMPATIBILITY.md](../COMPATIBILITY.md) 참조

## 라이선스

MIT — [LICENSE](../LICENSE).
