# 기여 가이드 — Gitea 메일 템플릿

## 기여 방법

### 새 스타일 추가하기

1. 도구로 스캐폴드: `cd tools && go run . create <style-name>` — 11개 이메일 유형 모두에 대한 플레이스홀더 `.tmpl` 파일이 포함된 완전한 디렉토리 구조를 생성합니다
2. `themes/<style-name>/`에서 각 `.tmpl` 파일을 고유한 디자인으로 편집
3. 프리뷰 재생성: `cd tools && go run . preview all` (빌드 스크립트가 `themes/` 아래의 모든 테마 디렉토리를 자동으로 찾고 테마 선택기도 자동 생성합니다)
4. 스크린샷과 함께 PR 제출 (이미지당 50 KiB 이하, 10–20 KiB 권장)

### 스타일 가이드라인

- 각 스타일에 **11개 템플릿 유형 모두** 포함
- Gitea 내장 함수만 사용
- 번역 키는 Gitea 공식 로케일 (`mail.*`) 사용
- **`.DisplayName` 사용 금지** (collaborator, transfer, release, workflow_run, assigned, default)
- 600px 너비의 이메일 클라이언트에 맞게 디자인
- Gmail, Outlook, Apple Mail에서 테스트

### 버그 리포트

1. Go 변수가 올바른지 확인
2. 번역 키가 Gitea 로케일과 일치하는지 확인
3. `.DisplayName`이 오용되지 않았는지 확인
4. 프리뷰 재생성: `cd tools && go run . preview all`
5. 스타일 이름과 메일 유형을 명시하여 이슈 생성

---

## 개발 설정

- **Go 1.21+** 템플릿 렌더링 및 CLI 도구
- **Node.js 18+** (선택) 개발 서버 및 Juice CSS 인라인

### 로컬 프리뷰 (정적)

1. 데이터 생성: `cd tools && go run . preview all`
2. `preview/index.html`을 브라우저에서 열기

> 정적 Gmail/Outlook 시뮬레이션은 참고용입니다. 정확한 렌더링은 dev 모드를 사용하세요.

### 개발 서버 (실시간 리로드 + CSS 인라인)

```bash
cd tools && go run . dev
# → http://localhost:3456
```

`.tmpl` 파일 수정 시 자동 재빌드되어 브라우저에 반영됩니다.

### 통합 테스트

Gitea에 배포: **Site Administration > Configuration > Mailer > Send Test Email**

---

## 커밋 규칙

- `style(<name>):` — 스타일 템플릿 변경
- `preview:` — 프리뷰 도구 변경
- `tools:` — Go 빌드 스크립트 변경
- `docs:` — 문서 및 번역
- `fix:` — 버그 수정
- `project:` — README, LICENSE, 메타 파일

## 라이선스

기여함으로써 MIT 라이선스에 동의하게 됩니다.
